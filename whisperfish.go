package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
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
	window        *qml.Window
	engine        *qml.Engine
	contactsModel Contacts
	configDir     string
	dataDir       string
	settings      *Settings
}

func main() {
	if err := qml.SailfishRun(APPNAME, "", VERSION, runGui); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Sailfish application failed")
	}
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
	w.contactsModel.Init()
	w.getPhoneNumber()
	w.getVerificationCode()
	w.getStoragePassword()
}

// Initializes qml context
func (w *Whisperfish) Init(engine *qml.Engine) {
	w.engine = engine
	w.engine.Translator(fmt.Sprintf("/usr/share/%s/qml/i18n", APPNAME))

	w.configDir = filepath.Join(w.engine.SailfishGetConfigLocation(), APPNAME)
	w.dataDir = w.engine.SailfishGetDataLocation()

	os.MkdirAll(w.configDir, 0700)
	os.MkdirAll(w.dataDir, 0700)

	settingsFile := filepath.Join(w.configDir, "settings.yml")
	w.settings = &Settings{}

	if err := w.settings.Load(settingsFile); err != nil {
		// write out default settings file
		if err = w.settings.Save(settingsFile); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to write out default settings file")
		}
	}

	// initialize model delegates
	w.engine.Context().SetVar("whisperfish", w)
	w.engine.Context().SetVar("contactsModel", &w.contactsModel)
}

// Returns the GO runtime version used for building the application
func (w *Whisperfish) RuntimeVersion() string {
	return runtime.Version()
}

// Returns the Whisperfish application version
func (w *Whisperfish) Version() string {
	return VERSION
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

// Get the current page status
func (w *Whisperfish) getCurrentPageStatus() int {
	return w.window.Root().ObjectByName("main").Object("currentPage").Int("status")
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
