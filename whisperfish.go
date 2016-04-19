package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/janimo/textsecure"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ttacon/libphonenumber"
	"gopkg.in/qml.v1"
)

const (
	VERSION                  = "0.1.1"
	APPNAME                  = "harbour-whisperfish"
	PAGE_STATUS_INACTIVE     = 0
	PAGE_STATUS_ACTIVATING   = 1
	PAGE_STATUS_ACTIVE       = 2
	PAGE_STATUS_DEACTIVATING = 3
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
	settings        *Settings
	config          *textsecure.Config
	db              *sqlx.DB
	activeSessionID int64
}

func main() {
	if err := qml.SailfishRun(APPNAME, "", VERSION, runGui); err != nil {
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

	_, err = db.Exec(SESSION_SCHEMA)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(MESSAGE_SCHEMA)
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
		GetLocalContacts:    getSailfishContacts,
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

	w.RefreshContacts()
	w.RefreshSessions()

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
	qml.Changed(&w.sessionModel, &w.sessionModel.Length)

	err := w.sessionModel.Refresh(w.db, &w.contactsModel)
	if err != nil && err != sql.ErrNoRows {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch sessions from database")
	}

	qml.Changed(&w.sessionModel, &w.sessionModel.Length)
}

// Set active session
func (w *Whisperfish) SetSession(sessionID int64) {
	w.activeSessionID = sessionID

	err := MarkSessionRead(w.db, sessionID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   w.activeSessionID,
		}).Error("Failed to mark session read")
	}
}

// Refresh conversation model
func (w *Whisperfish) RefreshConversation(sessionID int64) {
	w.messageModel.Length = 0
	qml.Changed(&w.messageModel, &w.messageModel.Length)

	err := w.messageModel.RefreshConversation(w.db, sessionID)
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
	w.engine.Translator(fmt.Sprintf("/usr/share/%s/qml/i18n", APPNAME))

	w.configDir = filepath.Join(w.engine.SailfishGetConfigLocation(), APPNAME)
	w.dataDir = w.engine.SailfishGetDataLocation()
	w.storageDir = filepath.Join(w.dataDir, "storage")
	dbDir := filepath.Join(w.dataDir, "db")
	dbFile := filepath.Join(dbDir, fmt.Sprintf("%s.db", APPNAME))

	os.MkdirAll(w.configDir, 0700)
	os.MkdirAll(w.dataDir, 0700)
	os.MkdirAll(dbDir, 0700)

	settingsFile := filepath.Join(w.configDir, "settings.yml")
	w.settings = &Settings{}

	if err := w.settings.Load(settingsFile); err != nil {
		w.settings.SetDefault()
		// write out default settings file
		if err = w.settings.Save(settingsFile); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to write out default settings file")
		}
	}

	var err error
	w.db, err = NewDb(dbFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to open database")
	}

	// initialize model delegates
	w.engine.Context().SetVar("whisperfish", w)
	w.engine.Context().SetVar("contactsModel", &w.contactsModel)
	w.engine.Context().SetVar("sessionModel", &w.sessionModel)
	w.engine.Context().SetVar("messageModel", &w.messageModel)
}

// Returns the GO runtime version used for building the application
func (w *Whisperfish) RuntimeVersion() string {
	return runtime.Version()
}

// Returns the Whisperfish application version
func (w *Whisperfish) Version() string {
	return VERSION
}

// Get the config file for Signal
func (w *Whisperfish) getConfig() (*textsecure.Config, error) {
	w.configFile = filepath.Join(w.configDir, "config.yml")
	var errConfig error
	if _, err := os.Stat(w.configFile); err == nil {
		w.config, errConfig = textsecure.ReadConfig(w.configFile)
	} else {
		w.config = &textsecure.Config{}
	}

	w.config.StorageDir = w.storageDir
	w.config.UserAgent = fmt.Sprintf("Whisperfish v%s", VERSION)
	w.config.UnencryptedStorage = true
	w.config.LogLevel = "debug"
	w.config.AlwaysTrustPeerID = true
	rootCA := filepath.Join(w.configDir, "rootCA.crt")
	if _, err := os.Stat(rootCA); err == nil {
		w.config.RootCA = rootCA
	}
	return w.config, errConfig
}

// Prompt the user for storage password
func (w *Whisperfish) getStoragePassword() string {
	pass := w.getTextFromDialog("getStoragePassword", "passwordDialog", "passwordEntered")
	log.Printf("Password: %s", pass)

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

	log.Println("Registered")
	status := w.getCurrentPageStatus()
	for status == PAGE_STATUS_ACTIVATING || status == PAGE_STATUS_DEACTIVATING {
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

// Get text from dialog window
func (w *Whisperfish) getTextFromDialog(fun, obj, signal string) string {
	status := w.getCurrentPageStatus()
	for status == PAGE_STATUS_ACTIVATING || status == PAGE_STATUS_DEACTIVATING {
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
	log.Printf("Recieved message from: %s", msg.Source())

	session, err := w.sessionModel.Add(w.db, msg.Source(), msg.Message(), time.Now(), true, false, false)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch session from database")
		return
	}

	message := &Message{
		SID:       session.ID,
		Source:    msg.Source(),
		Message:   msg.Message(),
		Timestamp: time.Now(),
	}

	err = SaveMessage(w.db, message)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   session.ID,
		}).Error("Failed to save message to database")
		return
	}

	if w.activeSessionID == session.ID && w.getCurrentPageID() == "conversation" {
		w.RefreshConversation(session.ID)
		err := MarkSessionRead(w.db, session.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"sid":   w.activeSessionID,
			}).Error("Failed to mark session read")
		}
	}

	w.RefreshSessions()

	if w.settings.EnableNotify {
		w.notify(msg)
	}
}

// Send new message notification
// From https://lists.sailfishos.org/pipermail/devel/2016-April/007036.html
func (w *Whisperfish) notify(msg *textsecure.Message) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	title := w.contactsModel.Name(msg.Source())
	body := "New message"

	var m map[string]dbus.Variant
	m = make(map[string]dbus.Variant)
	m["category"] = dbus.MakeVariant("x-nemo.messaging.im")

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "",
		uint32(0),
		"", title, body, []string{},
		m,
		int32(0))
	if call.Err != nil {
		return err
	}
	return nil
}

// Receipt handler
func (w *Whisperfish) receiptHandler(source string, devID uint32, timestamp uint64) {
	log.Printf("Receipt handler source %s timestamp %d", source, timestamp)
}
