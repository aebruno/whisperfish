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

package settings

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v2"
)

const (
	WhisperSettings = "harbour-whisperfish.conf"
)

//go:generate qtmoc
type Settings struct {
	core.QObject

	_ func(key string, val string)   `slot:"stringSet"`
	_ func(key string, val []string) `slot:"stringListSet"`
	_ func(key string, val bool)     `slot:"boolSet"`
	_ func(key string) string        `slot:"stringValue"`
	_ func(key string) []string      `slot:"stringListValue"`
	_ func(key string) bool          `slot:"boolValue"`
	_ func()                         `slot:"defaults"`

	_ func() `constructor:"init"`
}

// Initialize settings and setup slot connections
func (s *Settings) init() {
	// Connect slots
	s.ConnectStringSet(func(key string, val string) {
		s.SetString(key, val)
	})
	s.ConnectStringListSet(func(key string, val []string) {
		s.SetStringList(key, val)
	})
	s.ConnectBoolSet(func(key string, val bool) {
		s.SetBool(key, val)
	})
	s.ConnectStringValue(func(key string) string {
		return s.GetString(key)
	})
	s.ConnectStringListValue(func(key string) []string {
		return s.GetStringList(key)
	})
	s.ConnectBoolValue(func(key string) bool {
		return s.GetBool(key)
	})
	s.ConnectDefaults(func() {
		s.SetDefaults()
	})
}

// Set string value
func (s *Settings) SetString(key string, val string) {
	var settings = core.NewQSettings5(nil)
	settings.SetValue(key, core.NewQVariant14(val))
}

// Set string list value
func (s *Settings) SetStringList(key string, val []string) {
	var settings = core.NewQSettings5(nil)
	settings.SetValue(key, core.NewQVariant19(val))
}

// Set bool value
func (s *Settings) SetBool(key string, val bool) {
	var settings = core.NewQSettings5(nil)
	settings.SetValue(key, core.NewQVariant11(val))
}

// Get string value
func (s *Settings) GetString(key string) string {
	var settings = core.NewQSettings5(nil)
	return settings.Value(key, core.NewQVariant14("")).ToString()
}

// Get string list value
func (s *Settings) GetStringList(key string) []string {
	var settings = core.NewQSettings5(nil)
	return settings.Value(key, core.NewQVariant19([]string{})).ToStringList()
}

// Get bool value
func (s *Settings) GetBool(key string) bool {
	var settings = core.NewQSettings5(nil)
	return settings.Value(key, core.NewQVariant11(false)).ToBool()
}

// Set default values
func (s *Settings) SetDefaults() {
	var settings = core.NewQSettings5(nil)
	settings.SetValue("incognito", core.NewQVariant11(false))
	settings.SetValue("enable_notify", core.NewQVariant11(true))
	settings.SetValue("show_notify_message", core.NewQVariant11(false))
	settings.SetValue("encrypt_database", core.NewQVariant11(true))
	settings.SetValue("save_attachments", core.NewQVariant11(true))
	settings.SetValue("share_contacts", core.NewQVariant11(true))
	settings.SetValue("country_code", core.NewQVariant14(""))
	settings.Sync()
}

func (s *Settings) migrateSettings(file string) error {

	type deprecatedSettings struct {
		Incognito         bool   `yaml:"incognito"`
		EnableNotify      bool   `yaml:"enable_notify"`
		ShowNotifyMessage bool   `yaml:"show_notify_message"`
		EncryptDatabase   bool   `yaml:"encrypt_database"`
		SaveAttachments   bool   `yaml:"save_attachments"`
		CountryCode       string `yaml:"country_code"`
	}

	ds := &deprecatedSettings{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, ds)
	if err != nil {
		return err
	}

	var settings = core.NewQSettings5(nil)
	settings.SetValue("incognito", core.NewQVariant11(ds.Incognito))
	settings.SetValue("enable_notify", core.NewQVariant11(ds.EnableNotify))
	settings.SetValue("share_contacts", core.NewQVariant11(true))
	settings.SetValue("show_notify_message", core.NewQVariant11(ds.ShowNotifyMessage))
	settings.SetValue("encrypt_database", core.NewQVariant11(ds.EncryptDatabase))
	settings.SetValue("save_attachments", core.NewQVariant11(ds.SaveAttachments))
	settings.SetValue("country_code", core.NewQVariant14(ds.CountryCode))

	settings.Sync()

	return nil
}

func (s *Settings) Sync() {
	var settings = core.NewQSettings5(nil)
	settings.Sync()
}

// Initialize whisperfish settings file. Migrates old config if found.
func (s *Settings) Setup(configDir, storageDir string) error {
	deprecatedConfig := filepath.Join(configDir, "settings.yml")
	if _, err := os.Stat(deprecatedConfig); err == nil {
		log.Info("Deprecated settings.yml file found. Converting to new settings")

		err := s.migrateSettings(deprecatedConfig)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"file":  deprecatedConfig,
			}).Warn("Failed to convert old settings file")
		} else {
			log.Info("Succesfully converted settings. Removing old settings.yml file")
			os.Remove(deprecatedConfig)
		}
	}

	if _, err := os.Stat(filepath.Join(configDir, WhisperSettings)); os.IsNotExist(err) {
		log.Info("Config file not found. Setting default values")
		s.SetDefaults()
	}

	attachDir := s.GetString("attachment_dir")
	if attachDir == "" {
		attachDir = filepath.Join(storageDir, "attachments")
		os.MkdirAll(attachDir, 0700)
		s.SetString("attachment_dir", attachDir)
	}

	stat, err := os.Stat(attachDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("Attachment dir does not exist: %s", err)
	} else if err != nil {
		return fmt.Errorf("Failed to read attachment dir: %s", err)
	} else if !stat.IsDir() {
		return fmt.Errorf("Invalid setting for attachment_dir. Path is not a directory")
	} else if unix.Access(attachDir, unix.W_OK) != nil {
		return fmt.Errorf("Invalid setting for attachment_dir. Directory is not writable")
	}

	log.Infof("Attachments dir: %s", attachDir)

	searchPaths := s.GetStringList("search_paths")
	if len(searchPaths) == 0 || (len(searchPaths) == 1 && searchPaths[0] == "") {
		log.Infof("Empty search path found. Using defaults")
		searchPaths = []string{
			filepath.Join(core.QStandardPaths_WritableLocation(core.QStandardPaths__HomeLocation), "Pictures"),
		}
		s.SetStringList("search_paths", searchPaths)
	}

	for _, p := range searchPaths {
		log.Infof("Checking search path: %s", p)
		stat, err := os.Stat(p)
		if os.IsNotExist(err) {
			return fmt.Errorf("Search path does not exist: %s", err)
		} else if err != nil {
			return fmt.Errorf("Failed to read search path: %s", err)
		} else if !stat.IsDir() {
			return fmt.Errorf("Invalid setting for search_path. Path is not a directory")
		}
	}
	log.Infof("Attachment search paths: %s", searchPaths)
	return nil
}
