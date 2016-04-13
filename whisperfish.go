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
	Root qml.Object
}

var engine *qml.Engine
var win *qml.Window

func main() {
	if err := qml.SailfishRun(APPNAME, "", VERSION, run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Sailfish application failed")
	}
}

func runBackend() {
	refreshContacts()
}

func run() error {
	whisperfish := Whisperfish{}
	engine = qml.SailfishNewEngine()

	log.WithFields(log.Fields{
		"path": engine.SailfishGetConfigLocation(),
	}).Info("Configuration file location")
	log.WithFields(log.Fields{
		"path": engine.SailfishGetConfigLocation(),
	}).Info("Data file location")

	initModels()
	engine.Context().SetVar("whisperfish", &whisperfish)
	controls, err := engine.SailfishSetSource("qml/harbour-whisperfish.qml")
	if err != nil {
		return err
	}

	window := controls.SailfishCreateWindow()
	win = window
	whisperfish.Root = window.Root()

	window.SailfishShow()

	go runBackend()

	window.Wait()

	return nil
}

// Returns the GO runtime version used for building the application
func (w *Whisperfish) RuntimeVersion() string {
	return runtime.Version()
}

// Returns the Whisperfish application version
func (w *Whisperfish) Version() string {
	return VERSION
}
