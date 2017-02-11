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

package ui

import (
	"os"
	"path/filepath"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/model"
	"github.com/aebruno/whisperfish/settings"
	"github.com/aebruno/whisperfish/store"
	"github.com/aebruno/whisperfish/worker"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/sailfish"
)

const (
	WhisperDB   = "harbour-whisperfish.db"
	WhisperSalt = "salt"
)

// Create new DataStore
func NewStorage(dataPath, password string) (*store.DataStore, error) {
	// Set more restrictive umask to ensure database files are created 0600
	syscall.Umask(0077)

	var settings = settings.NewSettings(nil)

	dbDir := filepath.Join(dataPath, "db")
	log.Infof("Database dir: %s", dbDir)
	os.MkdirAll(dbDir, 0700)

	dbFile := filepath.Join(dbDir, WhisperDB)
	saltFile := ""

	if password != "" {
		saltFile = filepath.Join(dbDir, WhisperSalt)
	}

	if settings.GetBool("incognito") {
		dbFile = ":memory:"
	}

	ds, err := store.NewDataStore(dbFile, saltFile, password)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func InitDirs(configPath, dataPath string) (string, string) {
	log.Info("Setting up whisperfish directories")
	configDir := filepath.Join(configPath, "harbour-whisperfish")
	storageDir := filepath.Join(dataPath, "storage")

	log.Infof("Config dir: %s", configDir)
	os.MkdirAll(configDir, 0700)

	log.Infof("Data dir: %s", dataPath)
	os.MkdirAll(dataPath, 0700)

	return configDir, storageDir
}

// Run the main QT gui thread
func Run(version string) {
	log.Infof("Starting Whisperfish version %s", version)
	app := sailfish.SailfishApp_Application(len(os.Args), os.Args)
	app.SetApplicationVersion(version)
	app.SetOrganizationName("")
	app.SetApplicationName("harbour-whisperfish")

	// Setup i18n Translations
	var translator = core.NewQTranslator(nil)
	if translator.Load2(core.QLocale_System(), "whisperfish", "_", "/usr/share/harbour-whisperfish/qml/i18n", ".qm") {
		core.QCoreApplication_InstallTranslator(translator)
		log.WithFields(log.Fields{
			"locale": core.QLocale_System().Name(),
		}).Info("Successfully loaded system locale")
	} else {
		translator.Load2(core.NewQLocale3(core.QLocale__English, core.QLocale__UnitedStates), "whisperfish", "_", "/usr/share/harbour-whisperfish/qml/i18n", ".qm")
		core.QCoreApplication_InstallTranslator(translator)
		log.WithFields(log.Fields{
			"locale": core.QLocale_System().Name(),
		}).Info("No translations found for system locale. Using default")
	}

	var view = sailfish.SailfishApp_CreateView()

	var configPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__ConfigLocation)
	var dataPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__DataLocation)

	configDir, storageDir := InitDirs(configPath, dataPath)

	var settings = settings.NewSettings(nil)
	err := settings.Setup(configDir, storageDir)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to initialize settings")
	}

	var setupWorker = worker.NewSetupWorker(nil)
	config, err := setupWorker.ParseConfig(configDir, storageDir)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to parse signal config")
	}

	var filePicker = model.NewFilePicker(nil)
	var contactModel = model.NewContact(nil)
	var prompt = model.NewPrompt(nil)
	var sessionModel = model.NewSessionModel(nil)
	var messageModel = model.NewMessageModel(nil)
	var deviceModel = model.NewDeviceModel(nil)
	var clientWorker = worker.NewClientWorker(nil)
	var sendWorker = worker.NewSendWorker(nil)

	sendWorker.SetConfig(config)

	// Implement proper copy to clipboard support
	messageModel.ConnectCopyToClipboard(func(text string) {
		log.Info("Copy to clipboard")
		if len(text) > 0 {
			app.Clipboard().Clear(gui.QClipboard__Clipboard)
			app.Clipboard().SetText(text, gui.QClipboard__Clipboard)
		}
	})

	// Setup context variables
	view.RootContext().SetContextProperty("Prompt", prompt)
	view.RootContext().SetContextProperty("SettingsBridge", settings)
	view.RootContext().SetContextProperty("FilePicker", filePicker)
	view.RootContext().SetContextProperty("SessionModel", sessionModel)
	view.RootContext().SetContextProperty("MessageModel", messageModel)
	view.RootContext().SetContextProperty("ContactModel", contactModel)
	view.RootContext().SetContextProperty("DeviceModel", deviceModel)
	view.RootContext().SetContextProperty("SetupWorker", setupWorker)
	view.RootContext().SetContextProperty("ClientWorker", clientWorker)
	view.RootContext().SetContextProperty("SendWorker", sendWorker)

	client := &textsecure.Client{
		GetConfig: func() (*textsecure.Config, error) {
			return config, nil
		},
		GetPhoneNumber: func() string {
			return prompt.GetPhoneNumber()
		},
		GetVerificationCode: func() string {
			return prompt.GetVerificationCode()
		},
		GetStoragePassword: func() string {
			password := prompt.GetStoragePassword()

			if settings.GetBool("encrypt_database") {
				log.Info("Attempting to open encrypted datastore")
				var err error
				store.DS, err = NewStorage(dataPath, password)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Error("Failed to open encrypted database")
				}
			}

			return password
		},
		MessageHandler: func(msg *textsecure.Message) {
			log.Debug("Message received handler")
			clientWorker.MessageHandler(msg, false, 0)
		},
		ReceiptHandler: func(source string, devID uint32, ts uint64) {
			log.Debug("Message receipt handler")
			clientWorker.ReceiptHandler(source, devID, ts)
		},
		SyncReadHandler: func(source string, ts uint64) {
			log.Debug("Sync read handler")
			// TODO: not sure this is correct?
			clientWorker.ReceiptHandler(source, 0, ts)
		},
		SyncSentHandler: func(msg *textsecure.Message, ts uint64) {
			log.Debug("Sync sent handler")
			clientWorker.MessageHandler(msg, true, ts)
		},
	}

	setupWorker.ConnectSetupComplete(func() {
		if store.DS == nil && !settings.GetBool("encrypt_database") {
			log.Info("Attempting to open unencrypted datastore")
			var err error
			store.DS, err = NewStorage(dataPath, "")
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed to open unencrypted database")
			}
		}

		if store.DS == nil {
			log.Error("No datastore found")
			setupWorker.InvalidDatastore()
			return
		}

		setupWorker.SetLocked(false)
		contactModel.Refresh()
		sessionModel.Reload()
		deviceModel.Reload()
		clientWorker.Reconnect()
	})

	messageModel.ConnectSendMessage(func(mid int64) {
		// Only try and send messages when a network connection is established
		if clientWorker.IsConnected() {
			sendWorker.SendMessage(mid)
		}
	})

	view.SetSource(sailfish.SailfishApp_PathTo("qml/harbour-whisperfish.qml"))
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.Show()

	go setupWorker.Run(client)

	gui.QGuiApplication_Exec()
}
