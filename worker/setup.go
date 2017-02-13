// Copyright 2016 Andrew E. Bruno
//
// This file is part of Whisperfish.
//
// Whisperfish is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Whisperfish is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Whisperfish.  If not, see <http://www.gnu.org/licenses/>.

package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/settings"
	"github.com/aebruno/whisperfish/store"
	"github.com/therecipe/qt/core"
	"github.com/ttacon/libphonenumber"
)

const (
	SignalConfig = "config.yml"
)

//go:generate qtmoc
type SetupWorker struct {
	core.QObject

	settings   *settings.Settings
	config     *textsecure.Config
	configFile string

	_ bool   `property:"locked"`
	_ bool   `property:"registered"`
	_ string `property:"phoneNumber"`
	_ string `property:"localId"`
	_ string `property:"identity"`
	_ bool   `property:"encryptedKeystore"`
	_ func() `constructor:"init"`
	_ func() `signal:"registrationSuccess"`
	_ func() `signal:"setupComplete"`
	_ func() `signal:"invalidPhoneNumber"`
	_ func() `signal:"invalidDatastore"`
	_ func() `signal:"clientFailed"`
	_ func() `slot:"restart"`
}

// Setup connections
func (s *SetupWorker) init() {
	s.SetLocked(true)
	s.SetRegistered(false)
	s.SetEncryptedKeystore(false)
	s.settings = settings.NewSettings(nil)

	s.ConnectRestart(func() {
		s.settings.Sync()
		os.Exit(0)
	})
}

// Parse Signal config and create if not found
func (s *SetupWorker) ParseConfig(configDir, storageDir string) (*textsecure.Config, error) {
	log.Info("Parsing Signal config")
	log.Infof("Storage dir: %s", storageDir)
	os.MkdirAll(storageDir, 0700)

	s.config = &textsecure.Config{}
	s.configFile = filepath.Join(configDir, SignalConfig)

	if _, err := os.Stat(s.configFile); err == nil {
		s.config, err = textsecure.ReadConfig(s.configFile)
		if err != nil {
			return nil, err
		}
	} else {
		// Set defaults
		s.config.StorageDir = storageDir
		s.config.UserAgent = fmt.Sprintf("Whisperfish")
		s.config.UnencryptedStorage = false
		s.config.VerificationType = "voice"
		s.config.LogLevel = "debug"
		s.config.AlwaysTrustPeerID = false
	}

	rootCA := filepath.Join(configDir, "rootCA.crt")
	if _, err := os.Stat(rootCA); err == nil {
		s.config.RootCA = rootCA
	}

	log.Infof("Server: %s", s.config.Server)

	if _, err := os.Stat(filepath.Join(storageDir, "identity", "identity_key")); err == nil {
		log.Infof("Identity key found. Already registered")
		s.SetRegistered(true)
	}

	return s.config, nil
}

// Run Signal client setup
func (s *SetupWorker) Run(client *textsecure.Client) {
	log.Info("Setting up whisperfish client")

	client.GetConfig = func() (*textsecure.Config, error) {
		return s.config, nil
	}

	client.RegistrationDone = func() {
		log.Debug("Registration done handler")
		s.registrationDone()
	}
	client.GetLocalContacts = func() ([]textsecure.Contact, error) {
		log.Debug("Get local contacts handler")
		if s.settings.GetBool("share_contacts") {
			return store.SailfishContacts(s.settings.GetString("country_code"))
		}

		return make([]textsecure.Contact, 0), nil
	}

	err := textsecure.Setup(client)
	if err != nil {
		if _, ok := err.(*strconv.NumError); ok {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Invalid phone number in config file. Re-registration with Signal is required")
			s.InvalidPhoneNumber()
			return
		}

		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to setup Signal client")
		s.ClientFailed()
		return
	}

	s.SetPhoneNumber(s.phoneNumber())
	s.SetLocalId(s.config.Tel)
	s.SetIdentity(s.identity())
	s.SetEncryptedKeystore(!s.config.UnencryptedStorage)
	s.SetupComplete()
}

// Registration handler
func (s *SetupWorker) registrationDone() {
	textsecure.WriteConfig(s.configFile, s.config)

	num, err := libphonenumber.Parse(s.config.Tel, "")
	if err == nil {
		s.settings.SetString("country_code", libphonenumber.GetRegionCodeForNumber(num))
	}

	log.Info("Registered")
	s.SetRegistered(true)
	s.RegistrationSuccess()
}

// Returns the registered phone number
func (s *SetupWorker) phoneNumber() string {
	if s.config == nil {
		return ""
	}

	num, err := libphonenumber.Parse(s.config.Tel, "")
	if err == nil {
		return libphonenumber.Format(num, libphonenumber.NATIONAL)
	}

	return s.config.Tel
}

// Returns identity
func (s *SetupWorker) identity() string {
	id := textsecure.MyIdentityKey()
	return fmt.Sprintf("% 0X", id)
}
