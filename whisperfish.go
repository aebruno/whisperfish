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

package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/janimo/textsecure"
	"github.com/janimo/textsecure/3rd_party/magic"
	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/ttacon/libphonenumber"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/qml.v1"
)

const (
	Version                = "0.1.1"
	Appname                = "harbour-whisperfish"
	PageStatusInactive     = 0
	PageStatusActivating   = 1
	PageStatusActive       = 2
	PageStatusDeactivating = 3
	QtcontactsPath         = "/home/nemo/.local/share/system/Contacts/qtcontacts-sqlite/contacts.db"
)

type Whisperfish struct {
	window          *qml.Window
	engine          *qml.Engine
	contactsModel   Contacts
	sessionModel    SessionModel
	messageModel    MessageModel
	configDir       string
	configFile      string
	dataDir         string
	storageDir      string
	attachDir       string
	settingsFile    string
	settings        *Settings
	config          *textsecure.Config
	dbFile          string
	saltFile        string
	db              *sqlx.DB
	activeSessionID int64
	sentQueueSize   int
	totalMessages   int
	hasKeys         bool
	Locked          bool
}

func main() {
	if err := qml.SailfishRun(Appname, "", Version, runGui); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Sailfish application failed")
	}
}

func NewDb(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(SessionSchema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(MessageSchema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(SentqSchema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runGui() error {
	whisperfish := Whisperfish{}
	whisperfish.Init(qml.SailfishNewEngine())

	controls, err := whisperfish.engine.SailfishSetSource("qml/harbour-whisperfish.qml")
	if err != nil {
		return err
	}

	window := controls.SailfishCreateWindow()
	whisperfish.window = window

	window.SailfishShow()

	go whisperfish.runBackend()

	window.Wait()

	return nil
}

// Runs backend
func (w *Whisperfish) runBackend() {
	client := &textsecure.Client{
		GetConfig:           func() (*textsecure.Config, error) { return w.getConfig() },
		GetPhoneNumber:      func() string { return w.getPhoneNumber() },
		GetVerificationCode: func() string { return w.getVerificationCode() },
		GetStoragePassword:  func() string { return w.getStoragePassword() },
		MessageHandler:      func(msg *textsecure.Message) { w.messageHandler(msg) },
		ReceiptHandler:      func(source string, devID uint32, timestamp uint64) { w.receiptHandler(source, devID, timestamp) },
		RegistrationDone:    func() { w.registrationDone() },
		GetLocalContacts:    func() ([]textsecure.Contact, error) { return w.getSailfishContacts() },
	}

	err := textsecure.Setup(client)
	if _, ok := err.(*strconv.NumError); ok {
		os.RemoveAll(w.storageDir)
		log.Fatal("Switching to unencrypted session store, removing %s\nThis will reset your sessions and reregister your phone.", w.storageDir)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to setup textsecure client")
		return
	}

	if w.db == nil && !w.settings.EncryptDatabase {
		// Attempt open of unencrypted datastore
		var err error
		w.db, err = NewDb(w.dbFile)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to open unencrypted database")
			return
		}
	}

	if w.db == nil {
		log.Error("No database found")
		return
	}

	w.Locked = false
	qml.Changed(w, &w.Locked)
	w.RefreshContacts()
	w.RefreshSessions()

	go w.sendMessageWorker()

	for {
		if err := textsecure.StartListening(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error processing Websocket event from Signal")
			time.Sleep(3 * time.Second)
		}
	}
}

// Refresh contacts
func (w *Whisperfish) RefreshContacts() {
	w.contactsModel.Refresh()
}

// Refresh session model
func (w *Whisperfish) RefreshSessions() {
	w.sessionModel.Length = 0
	w.sessionModel.Unread = 0
	qml.Changed(&w.sessionModel, &w.sessionModel.Length)
	qml.Changed(&w.sessionModel, &w.sessionModel.Unread)

	err := w.sessionModel.Refresh(w.db, &w.contactsModel)
	if err != nil && err != sql.ErrNoRows {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch sessions from database")
	}

	qml.Changed(&w.sessionModel, &w.sessionModel.Length)
	qml.Changed(&w.sessionModel, &w.sessionModel.Unread)
}

// Set active session
func (w *Whisperfish) setSession(session *Session) {
	w.activeSessionID = session.ID
	if session.IsGroup {
		w.messageModel.Name = session.GroupName
		w.messageModel.Identity = ""
	} else {
		w.messageModel.Name = w.contactsModel.Name(session.Source)
		w.messageModel.Identity = w.ContactIdentity(session.Source)
	}
	w.messageModel.IsGroup = session.IsGroup
	w.messageModel.Tel = session.Source
	qml.Changed(&w.messageModel, &w.messageModel.Name)
	qml.Changed(&w.messageModel, &w.messageModel.IsGroup)
	qml.Changed(&w.messageModel, &w.messageModel.Tel)
	qml.Changed(&w.messageModel, &w.messageModel.Identity)
}

// Set active session by id
func (w *Whisperfish) SetSession(sessionID int64) {
	session, err := FetchSession(w.db, sessionID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    sessionID,
		}).Error("Failed to fetch session")
	}

	w.setSession(session)

	err = MarkSessionRead(w.db, session.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   session.ID,
		}).Error("Failed to mark session read")
	}

	w.messageModel.Length = 0
	qml.Changed(&w.messageModel, &w.messageModel.Length)
	w.RefreshConversation()
	w.RefreshSessions()
}

// Refresh conversation model
func (w *Whisperfish) RefreshConversation() {
	err := w.messageModel.RefreshConversation(w.db, w.activeSessionID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch messages from database")
	}

	qml.Changed(&w.messageModel, &w.messageModel.Length)
}

// Initializes Whisperfish application and qml context
func (w *Whisperfish) Init(engine *qml.Engine) {
	w.engine = engine
	w.engine.Translator(fmt.Sprintf("/usr/share/%s/qml/i18n", Appname))

	w.configDir = filepath.Join(w.engine.SailfishGetConfigLocation(), Appname)
	w.dataDir = w.engine.SailfishGetDataLocation()
	w.storageDir = filepath.Join(w.dataDir, "storage")
	w.attachDir = filepath.Join(w.storageDir, "attachments")
	dbDir := filepath.Join(w.dataDir, "db")
	w.dbFile = filepath.Join(dbDir, fmt.Sprintf("%s.db", Appname))
	w.saltFile = filepath.Join(dbDir, "salt")

	os.MkdirAll(w.configDir, 0700)
	os.MkdirAll(w.dataDir, 0700)
	os.MkdirAll(w.attachDir, 0700)
	os.MkdirAll(dbDir, 0700)

	w.settingsFile = filepath.Join(w.configDir, "settings.yml")
	w.settings = &Settings{}

	if err := w.settings.Load(w.settingsFile); err != nil {
		w.settings.SetDefault()
		w.SaveSettings()
	}

	if w.settings.Incognito {
		w.dbFile = ":memory:"
	}

	// initialize model delegates
	w.engine.Context().SetVar("whisperfish", w)
	w.engine.Context().SetVar("contactsModel", &w.contactsModel)
	w.engine.Context().SetVar("sessionModel", &w.sessionModel)
	w.engine.Context().SetVar("messageModel", &w.messageModel)

	// default to locked
	w.Locked = true

	if _, err := os.Stat(filepath.Join(w.storageDir, "identity", "identity_key")); err == nil {
		w.hasKeys = true
	}
}

// Force exit of application
func (w *Whisperfish) Restart() {
	os.Exit(0)
}

// Returns the GO runtime version used for building the application
func (w *Whisperfish) RuntimeVersion() string {
	return runtime.Version()
}

// Returns the Whisperfish application version
func (w *Whisperfish) Version() string {
	return Version
}

// Returns the registered phone number
func (w *Whisperfish) PhoneNumber() string {
	if w.config == nil {
		return ""
	}

	num, err := libphonenumber.Parse(w.config.Tel, "")
	if err == nil {
		return libphonenumber.Format(num, libphonenumber.NATIONAL)
	}

	return w.config.Tel
}

// Returns identity
func (w *Whisperfish) Identity() string {
	id := textsecure.MyIdentityKey()
	return fmt.Sprintf("% 0X", id)
}

// Returns identity of contact
func (w *Whisperfish) ContactIdentity(source string) string {
	id, err := textsecure.ContactIdentityKey(source)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch contact identity")
		return ""
	}

	return fmt.Sprintf("% 0X", id)
}

// Return settings
func (w *Whisperfish) Settings() *Settings {
	return w.settings
}

// Save settings
func (w *Whisperfish) SaveSettings() {
	if w.settings == nil {
		return
	}

	err := w.settings.Save(w.settingsFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to write settings file")
	}
}

// Delete message
func (w *Whisperfish) DeleteMessage(id int64) {
	err := DeleteMessage(w.db, id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to delete message")
	}
	DequeueSent(w.db, id)
	w.RefreshConversation()
}

// Get the config file for Signal
func (w *Whisperfish) getConfig() (*textsecure.Config, error) {
	w.configFile = filepath.Join(w.configDir, "config.yml")
	var errConfig error
	if _, err := os.Stat(w.configFile); err == nil {
		w.config, errConfig = textsecure.ReadConfig(w.configFile)
	} else {
		w.config = &textsecure.Config{}

		// Set defaults
		w.config.StorageDir = w.storageDir
		w.config.UserAgent = fmt.Sprintf("Whisperfish v%s", Version)
		w.config.UnencryptedStorage = false
		w.config.VerificationType = "voice"
		w.config.LogLevel = "debug"
		w.config.AlwaysTrustPeerID = true
	}

	rootCA := filepath.Join(w.configDir, "rootCA.crt")
	if _, err := os.Stat(rootCA); err == nil {
		w.config.RootCA = rootCA
	}
	return w.config, errConfig
}

func (w *Whisperfish) getSalt() ([]byte, error) {
	salt := make([]byte, 8)

	if _, err := os.Stat(w.saltFile); err == nil {
		salt, err = ioutil.ReadFile(w.saltFile)
		if err != nil {
			return nil, err
		}
	} else {
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(w.saltFile, salt, 0600)
		if err != nil {
			return nil, err
		}
	}

	return salt, nil
}

// Prompt the user for storage password
func (w *Whisperfish) getStoragePassword() string {
	pass := w.getTextFromDialog("getStoragePassword", "passwordDialog", "passwordEntered")

	if w.settings.EncryptDatabase {
		salt, err := w.getSalt()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to get salt")
			return pass
		}

		key, _ := scrypt.Key([]byte(pass), salt, 16384, 8, 1, 32)
		dsn := fmt.Sprintf("%s?_pragma_key=x'%X'&_pragma_cipher_page_size=4096", w.dbFile, key)
		w.db, err = NewDb(dsn)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to open encrypted database")
		}
	}

	return pass
}

// Prompt the user to enter the verification code
func (w *Whisperfish) getVerificationCode() string {
	code := w.getTextFromDialog("getVerificationCode", "verifyDialog", "codeEntered")
	log.Printf("Code: %s", code)

	return code
}

// Prompt the user to enter telephone number for Registration
func (w *Whisperfish) getPhoneNumber() string {
	n := w.getTextFromDialog("getPhoneNumber", "registerDialog", "numberEntered")
	num, err := libphonenumber.Parse(fmt.Sprintf("+%s", n), "")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to parse phone number")
	}

	tel := libphonenumber.Format(num, libphonenumber.E164)
	log.Printf("Using phone number: %s", tel)
	return tel
}

// Registration done
func (w *Whisperfish) registrationDone() {
	textsecure.WriteConfig(w.configFile, w.config)

	num, err := libphonenumber.Parse(w.config.Tel, "")
	if err == nil {
		w.settings.CountryCode = libphonenumber.GetRegionCodeForNumber(num)
		w.SaveSettings()
	}

	log.Println("Registered")
	status := w.getCurrentPageStatus()
	for status == PageStatusActivating || status == PageStatusDeactivating {
		// If current page is in transition need to wait before pushing dialog on stack
		time.Sleep(100 * time.Millisecond)
		status = w.getCurrentPageStatus()
	}
	w.window.Root().ObjectByName("main").Call("registered")
}

// Get the current page status
func (w *Whisperfish) getCurrentPageStatus() int {
	return w.window.Root().ObjectByName("main").Object("currentPage").Int("status")
}

// Get the current page id
func (w *Whisperfish) getCurrentPageID() string {
	return w.window.Root().ObjectByName("main").Object("currentPage").String("objectName")
}

// Returns true if applications is active
func (w *Whisperfish) isActive() bool {
	return w.window.Root().Bool("applicationActive")
}

// Get text from dialog window
func (w *Whisperfish) getTextFromDialog(fun, obj, signal string) string {
	status := w.getCurrentPageStatus()
	for status == PageStatusActivating || status == PageStatusDeactivating {
		// If current page is in transition need to wait before pushing dialog on stack
		time.Sleep(100 * time.Millisecond)
		status = w.getCurrentPageStatus()
	}

	w.window.Root().ObjectByName("main").Call(fun)
	p := w.window.Root().ObjectByName(obj)
	ch := make(chan string)
	p.On(signal, func(text string) {
		ch <- text
	})
	text := <-ch
	return text
}

// Message handler
func (w *Whisperfish) messageHandler(msg *textsecure.Message) {
	log.Printf("Received message from: %s", msg.Source())

	message := &Message{
		Source:    msg.Source(),
		Message:   msg.Message(),
		Timestamp: time.Now(),
		Flags:     msg.Flags(),
	}

	if len(msg.Attachments()) > 0 {
		if w.settings.SaveAttachments && !w.settings.Incognito {
			err := message.SaveAttachment(w.attachDir, msg.Attachments()[0])
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed to save attachment")
			}
		} else {
			message.HasAttachment = true
			message.MimeType = msg.Attachments()[0].MimeType
		}
	}

	session, err := w.sessionModel.Add(w.db, message, msg.Group(), true)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add message to database")
		return
	}

	if w.activeSessionID == session.ID {
		w.RefreshConversation()

		if w.getCurrentPageID() == "conversation" {
			err := MarkSessionRead(w.db, session.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"sid":   w.activeSessionID,
				}).Error("Failed to mark session read")
			}
		}
	}

	w.RefreshSessions()

	active := w.isActive()
	pageID := w.getCurrentPageID()

	// Don't send notification if disabled or viewing the main conversation page
	if !w.settings.EnableNotify || (active && pageID == "main") {
		return
	}

	// Don't send notification if view the current conversation
	if active && w.activeSessionID == session.ID && pageID == "conversation" {
		return
	}

	name := w.contactsModel.Name(msg.Source())
	w.window.Root().Call("newMessageNotification", session.ID, name, msg.Message())
}

// Send message
func (w *Whisperfish) SendMessage(source, message, groupName, attachment string) {
	var err error

	m := strings.Split(source, ",")
	if len(m) > 1 {
		group, err := textsecure.NewGroup(groupName, m)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"group_name": groupName,
			}).Error("Failed to create new group")
			return
		}

		err = w.sendMessageHelper(group.Hexid, message, attachment, group)
	} else {
		err = w.sendMessageHelper(source, message, attachment, nil)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   w.activeSessionID,
		}).Error("Failed to send message")
	}
}

func (w *Whisperfish) sendMessageHelper(to, msg, attachment string, group *textsecure.Group) error {
	message := &Message{
		Source:    to,
		Message:   msg,
		Timestamp: time.Now(),
		Outgoing:  true,
	}

	if len(attachment) > 0 {
		att, err := os.Open(attachment)
		if err != nil {
			return err
		}
		defer att.Close()
		//XXX Sucks we have to do this twice
		message.MimeType, _ = magic.MIMETypeFromReader(att)
		message.Attachment = attachment
		message.HasAttachment = true
	}

	session, err := w.sessionModel.Add(w.db, message, group, false)
	if err != nil {
		return err
	}

	w.setSession(session)
	w.RefreshConversation()
	w.RefreshSessions()

	err = QueueSent(w.db, message)
	if err != nil {
		return err
	}

	return nil
}

func (w *Whisperfish) sendMessage(m *Message) error {
	var att io.Reader
	var err error
	var ts uint64

	s, err := FetchSession(w.db, m.SID)
	if err != nil {
		return err
	}

	if m.Attachment != "" {
		att, err = os.Open(m.Attachment)
		if err != nil {
			return err
		}
	}

	if att == nil {
		if s.IsGroup {
			ts, err = textsecure.SendGroupMessage(s.Source, m.Message)
		} else {
			ts, err = textsecure.SendMessage(s.Source, m.Message)
		}
	} else {
		if s.IsGroup {
			ts, err = textsecure.SendGroupAttachment(s.Source, m.Message, att)
		} else {
			ts, err = textsecure.SendAttachment(s.Source, m.Message, att)
		}
	}

	if err != nil {
		return err
	}

	timestamp := time.Unix(0, int64(1000000*ts)).Local()

	err = MarkSessionSent(w.db, s.ID, m.Message, timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    s.ID,
		}).Error("Failed to mark session sent")
		return err
	}

	err = MarkMessageSent(w.db, m.ID, timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    m.ID,
		}).Error("Failed to mark message sent")
		return err
	}

	if w.activeSessionID == s.ID {
		for i, x := range w.messageModel.messages {
			if x.ID == m.ID {
				w.messageModel.messages[i].Sent = true
				qml.Changed(w.messageModel.messages[i], &w.messageModel.messages[i].Sent)
			}
		}
	}

	w.RefreshSessions()
	return nil
}

func (w *Whisperfish) sendMessageWorker() {
	for {
		messages, err := FetchSentq(w.db)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to fetch mailq")
		}

		w.sentQueueSize = len(messages)

		for _, m := range messages {
			err = w.sendMessage(m)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    m.ID,
				}).Error("Failed to send message")
				continue
			}

			// Remove from sentq
			err := DequeueSent(w.db, m.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    m.ID,
				}).Error("Failed to remove message from mailq")
			}

			// Throttle
			time.Sleep(1 * time.Second)
		}

		w.totalMessages, err = TotalMessages(w.db)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to update total messages")
		}

		time.Sleep(3 * time.Second)
	}
}

// Receipt handler
func (w *Whisperfish) receiptHandler(source string, devID uint32, ts uint64) {
	log.Printf("Receipt handler source %s timestamp %d", source, ts)

	timestamp := time.Unix(0, int64(1000000*ts)).Local()

	sessionID, messageID, err := MarkMessageReceived(w.db, source, timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to mark message received")
	}

	err = MarkSessionReceived(w.db, sessionID, timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    sessionID,
		}).Error("Failed to mark session received")
	}

	if w.activeSessionID == sessionID {
		for i, x := range w.messageModel.messages {
			if x.ID == messageID {
				w.messageModel.messages[i].Received = true
				qml.Changed(w.messageModel.messages[i], &w.messageModel.messages[i].Received)
			}
		}
	}

	w.RefreshSessions()
}

// Get local sailfish contacts
func (w *Whisperfish) getSailfishContacts() ([]textsecure.Contact, error) {
	db, err := sqlx.Open("sqlite3", QtcontactsPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to open contacts database")
		return nil, err
	}

	contacts := []textsecure.Contact{}
	err = db.Select(&contacts, `
	select
	   c.displayLabel as name,
	   p.phoneNumber as tel
	from Contacts as c
	join PhoneNumbers p
	on c.contactId = p.contactId`)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to query contacts database")
		return nil, err
	}

	// Reformat numbers in E.164 format
	for i := range contacts {
		n := contacts[i].Tel
		num, err := libphonenumber.Parse(n, w.settings.CountryCode)
		if err == nil {
			contacts[i].Tel = libphonenumber.Format(num, libphonenumber.E164)
		}
	}

	return contacts, nil
}

// Return true if encrypted keys have been created
func (w *Whisperfish) HasEncryptionKeys() bool {
	return w.hasKeys
}

// Return true if encrypted key store is enabled
func (w *Whisperfish) HasEncryptedKeystore() bool {
	if w.config == nil {
		return false
	}

	return !w.config.UnencryptedStorage
}

func (w *Whisperfish) TotalMessages() int {
	return w.totalMessages
}

func (w *Whisperfish) SentQueueSize() int {
	return w.sentQueueSize
}
