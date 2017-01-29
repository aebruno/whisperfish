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

package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/model"
	"github.com/aebruno/whisperfish/store"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/ttacon/libphonenumber"
	"golang.org/x/sys/unix"
)

const (
	SignalConfig    = "config.yml"
	WhisperSettings = "harbour-whisperfish.conf"
	WhisperDB       = "harbour-whisperfish.db"
	WhisperSalt     = "salt"
)

//go:generate qtmoc
type Backend struct {
	core.QObject

	settings     *model.Settings
	prompt       *model.Prompt
	sessionModel *model.SessionModel
	messageModel *model.MessageModel
	contacts     *store.Contacts
	config       *textsecure.Config
	dataDir      string
	configFile   string
	ds           *store.DataStore

	_ bool                                                `property:"locked"`
	_ bool                                                `property:"registered"`
	_ bool                                                `property:"connected"`
	_ func()                                              `constructor:"init"`
	_ func()                                              `signal:"registrationSuccess"`
	_ func(id int64, source, message string)              `signal:"notifyMessage"`
	_ func(text string)                                   `slot:"copyToClipboard"`
	_ func(tel string) string                             `slot:"contactNumber"`
	_ func(tel string) string                             `slot:"contactIdentity"`
	_ func(tel string) string                             `slot:"contactName"`
	_ func() int                                          `slot:"contactCount"`
	_ func()                                              `slot:"contactRefresh"`
	_ func() string                                       `slot:"phoneNumber"`
	_ func() string                                       `slot:"identity"`
	_ func() bool                                         `slot:"hasEncryptedKeystore"`
	_ func()                                              `slot:"reconnect"`
	_ func()                                              `slot:"restart"`
	_ func(source string)                                 `slot:"endSession"`
	_ func(source, message, groupName, attachment string) `slot:"sendMessage"`
	_ func(sid int64)                                     `slot:"activateSession"`
}

// Setup connections
func (b *Backend) init() {
	// Slot connections
	b.ConnectPhoneNumber(func() string {
		return b.phoneNumber()
	})
	b.ConnectIdentity(func() string {
		return b.identity()
	})
	b.ConnectHasEncryptedKeystore(func() bool {
		return b.hasEncryptedKeystore()
	})
	b.ConnectReconnect(func() {
		b.reconnect()
	})
	b.ConnectContactIdentity(func(source string) string {
		return b.contactIdentity(source)
	})
	b.ConnectContactNumber(func(tel string) string {
		return b.contacts.Find(tel, b.settings.GetString("country_code"))
	})
	b.ConnectContactName(func(tel string) string {
		return b.contacts.FindName(tel)
	})
	b.ConnectContactCount(func() int {
		return b.contacts.Len()
	})
	b.ConnectActivateSession(func(sid int64) {
		b.activateSession(sid)
	})
	b.ConnectContactRefresh(func() {
		b.contacts.Refresh()
		b.sessionModel.Load()
	})
	b.ConnectRestart(func() {
		b.settings.Sync()
		os.Exit(0)
	})
	b.ConnectEndSession(func(source string) {
		b.endSession(source)
	})
	b.ConnectSendMessage(func(source, message, groupName, attachment string) {
		b.sendMessage(source, message, groupName, attachment)
	})
}

// Setup config/settings directories, parse configs, and set context variables
func (b *Backend) Setup(configPath, dataPath string, view *quick.QQuickView) error {
	log.Info("Setting up backend")
	b.dataDir = dataPath
	configDir := filepath.Join(configPath, "harbour-whisperfish")
	storageDir := filepath.Join(dataPath, "storage")

	log.Infof("Config dir: %s", configDir)
	os.MkdirAll(configDir, 0700)

	log.Infof("Data dir: %s", dataPath)
	os.MkdirAll(dataPath, 0700)

	err := b.parseSignalConfig(configDir, storageDir)
	if err != nil {
		return err
	}

	err = b.initWhisperSettings(configDir, storageDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(storageDir, "identity", "identity_key")); err == nil {
		b.SetRegistered(true)
	}

	b.prompt = model.NewPrompt(nil)
	b.messageModel = model.NewMessageModel(nil)
	b.sessionModel = model.NewSessionModel(nil)

	if view != nil {
		var filePicker = model.NewFilePicker(nil)

		// Setup context variables
		view.RootContext().SetContextProperty("Backend", b)
		view.RootContext().SetContextProperty("Prompt", b.prompt)
		view.RootContext().SetContextProperty("SettingsBridge", b.settings)
		view.RootContext().SetContextProperty("FilePicker", filePicker)
		view.RootContext().SetContextProperty("FileModel", filePicker.Model)
		view.RootContext().SetContextProperty("SessionModel", b.sessionModel)
		view.RootContext().SetContextProperty("SessionListModel", b.sessionModel.Model)
		view.RootContext().SetContextProperty("MessageModel", b.messageModel)
		view.RootContext().SetContextProperty("MessageListModel", b.messageModel.Model)
	}

	return nil
}

// Parse Signal config and create if not found
func (b *Backend) parseSignalConfig(configDir, storageDir string) error {
	log.Info("Parsing Signal config")
	log.Infof("Storage dir: %s", storageDir)
	os.MkdirAll(storageDir, 0700)

	b.config = &textsecure.Config{}
	b.configFile = filepath.Join(configDir, SignalConfig)

	if _, err := os.Stat(b.configFile); err == nil {
		b.config, err = textsecure.ReadConfig(b.configFile)
		if err != nil {
			return err
		}
	} else {
		// Set defaults
		b.config.StorageDir = storageDir
		b.config.UserAgent = fmt.Sprintf("Whisperfish")
		b.config.UnencryptedStorage = false
		b.config.VerificationType = "voice"
		b.config.LogLevel = "debug"
		b.config.AlwaysTrustPeerID = false
	}

	rootCA := filepath.Join(configDir, "rootCA.crt")
	if _, err := os.Stat(rootCA); err == nil {
		b.config.RootCA = rootCA
	}

	log.Infof("Server: %s", b.config.Server)
	return nil
}

// Initialize whisperfish settings file. Uses QSettings. Migrates old config if
// found.
func (b *Backend) initWhisperSettings(configDir, storageDir string) error {
	b.settings = model.NewSettings(nil)

	deprecatedConfig := filepath.Join(configDir, "settings.yml")
	if _, err := os.Stat(deprecatedConfig); err == nil {
		log.Info("Deprecated settings.yml file found. Converting to new settings")

		err := b.settings.MigrateSettings(deprecatedConfig)
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
		b.settings.SetDefaults()
	}

	attachDir := b.settings.GetString("attachment_dir")
	if attachDir == "" {
		attachDir = filepath.Join(storageDir, "attachments")
		os.MkdirAll(attachDir, 0700)
	}

	stat, err := os.Stat(attachDir)
	if os.IsNotExist(err) {
		attachDir = filepath.Join(storageDir, "attachments")
		b.settings.SetString("attachment_dir", attachDir)
	} else if err != nil {
		return fmt.Errorf("Failed to read attachment dir: %s", err)
	} else if !stat.IsDir() {
		return fmt.Errorf("Invalid setting for attachment_dir. Path is not a directory")
	} else if unix.Access(attachDir, unix.W_OK) != nil {
		return fmt.Errorf("Invalid setting for attachment_dir. Directory is not writable")
	}

	log.Infof("Attachments dir: %s", attachDir)
	return nil
}

// Create new DataStore
func (b *Backend) newStorage(password string) error {
	dbDir := filepath.Join(b.dataDir, "db")
	log.Infof("Database dir: %s", dbDir)
	os.MkdirAll(dbDir, 0700)

	dbFile := filepath.Join(dbDir, WhisperDB)
	saltFile := ""

	if password != "" {
		saltFile = filepath.Join(dbDir, WhisperSalt)
	}

	var err error
	b.ds, err = store.NewDataStore(dbFile, saltFile, password)
	if err != nil {
		return err
	}

	return nil
}

// Run Signal client thread that connects view websockets
func (b *Backend) Run() {
	log.Info("Starting whisperfish backend thread")

	client := &textsecure.Client{
		GetConfig:           func() (*textsecure.Config, error) { return b.config, nil },
		GetPhoneNumber:      func() string { return b.prompt.GetPhoneNumber() },
		GetVerificationCode: func() string { return b.prompt.GetVerificationCode() },
		GetStoragePassword:  func() string { return b.getStoragePassword() },
		MessageHandler:      func(msg *textsecure.Message) { b.messageHandler(msg) },
		ReceiptHandler:      func(source string, devID uint32, timestamp uint64) { b.receiptHandler(source, devID, timestamp) },
		RegistrationDone:    func() { b.registrationDone() },
		GetLocalContacts:    func() ([]textsecure.Contact, error) { return b.getLocalContacts() },
		SyncReadHandler:     func(source string, ts uint64) { b.syncReadHandler(source, ts) },
		SyncSentHandler:     func(msg *textsecure.Message, ts uint64) { b.syncSentHandler(msg, ts) },
	}

	err := textsecure.Setup(client)
	if err != nil {
		if _, ok := err.(*strconv.NumError); ok {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Invalid phone number in config file. Re-registration with Signal is required")
		}

		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to setup textsecure client")
		return
	}

	if b.ds == nil && !b.settings.GetBool("encrypt_database") {
		log.Info("Attempting to open unencrypted datastore")
		err := b.newStorage("")
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to open unencrypted database")
			return
		}
	}

	if b.ds == nil {
		log.Error("No datastore found")
		return
	}

	b.SetLocked(false)
	b.contacts = store.NewContacts(b.settings.GetString("country"))
	b.sessionModel.SetDataStore(b.ds)
	b.messageModel.SetDataStore(b.ds)
	b.sessionModel.Refresh()

	go b.sendMessageWorker()

	for {
		time.Sleep(3 * time.Second)

		if !b.IsConnected() {
			log.Debug("No network connection found")
			continue
		}

		log.Debug("Starting textsecure websocket listener")
		if err := textsecure.StartListening(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error processing Websocket event from Signal")
		}
	}
}

// Return true if encrypted key store is enabled
func (b *Backend) hasEncryptedKeystore() bool {
	if b.config == nil {
		return false
	}

	return !b.config.UnencryptedStorage
}

// Return true if encrypted database is enabled
func (b *Backend) HasEncryptedDatabase() bool {
	return b.settings.GetBool("encrypt_database")
}

// Returns identity of contact
func (b *Backend) contactIdentity(source string) string {
	id, err := textsecure.ContactIdentityKey(source)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch contact identity")
		return ""
	}

	return fmt.Sprintf("% 0X", id)
}

// Force reconnect of websocket connection
func (b *Backend) reconnect() {
	log.Info("Forcing websocket reconnection")
	textsecure.StopListening()
}

// Returns the registered phone number
func (b *Backend) phoneNumber() string {
	if b.config == nil {
		return ""
	}

	num, err := libphonenumber.Parse(b.config.Tel, "")
	if err == nil {
		return libphonenumber.Format(num, libphonenumber.NATIONAL)
	}

	return b.config.Tel
}

// Returns identity
func (b *Backend) identity() string {
	id := textsecure.MyIdentityKey()
	return fmt.Sprintf("% 0X", id)
}

// Activate session by id
func (b *Backend) activateSession(sid int64) {
	s, err := b.ds.FetchSession(sid)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   sid,
		}).Error("No session found")
		return
	}

	b.sessionModel.MarkRead(sid)
	b.messageModel.Refresh(sid, b.contacts.Find(s.Source, b.settings.GetString("country_code")), b.contactIdentity(s.Source), s.Source, s.IsGroup)
}

// Convert data store
func (b *Backend) ConvertDataStore(password string) error {
	if password == "" {
		return fmt.Errorf("No password given")
	}

	dbDir := filepath.Join(b.dataDir, "db")
	tmp := filepath.Join(dbDir, "tmp.db")
	encrypt := !b.settings.GetBool("encrypt_database")

	if encrypt {
		log.Info("Encrypting database..")

		err := b.newStorage("")
		if err != nil {
			return err
		}

		err = b.ds.Encrypt(tmp, password)
		if err != nil {
			return err
		}
	} else {
		log.Info("Decrypting database..")

		err := b.newStorage(password)
		if err != nil {
			return err
		}

		err = b.ds.Decrypt(tmp)
		if err != nil {
			return err
		}
	}

	dbFile := filepath.Join(dbDir, WhisperDB)

	err := os.Rename(tmp, dbFile)
	if err != nil {
		return err
	}

	b.settings.SetBool("encrypt_database", encrypt)
	b.settings.Sync()
	return nil
}
