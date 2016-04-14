package main

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/qml.v1"
)

const (
	VERSION = "0.1.1"
	APPNAME = "harbour-whisperfish"
)

type Whisperfish struct {
	window        *qml.Window
	engine        *qml.Engine
	contactsModel Contacts
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

	log.WithFields(log.Fields{
		"path": whisperfish.engine.SailfishGetConfigLocation(),
	}).Info("Configuration file location")
	log.WithFields(log.Fields{
		"path": whisperfish.engine.SailfishGetDataLocation(),
	}).Info("Data file location")

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
}

// Initializes qml context
func (w *Whisperfish) Init(engine *qml.Engine) {
	w.engine = engine
	w.engine.Context().SetVar("whisperfish", w)
	w.engine.Context().SetVar("contactsModel", &w.contactsModel)
	w.engine.Translator(fmt.Sprintf("/usr/share/%s/qml/i18n", APPNAME))
}

// Returns the GO runtime version used for building the application
func (w *Whisperfish) RuntimeVersion() string {
	return runtime.Version()
}

// Returns the Whisperfish application version
func (w *Whisperfish) Version() string {
	return VERSION
}
