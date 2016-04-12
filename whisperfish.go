package main

import (
	"fmt"
	"gopkg.in/qml.v1"
	"os"
	"runtime"
)

const VERSION = "0.1.1"

type Whisperfish struct {
	Root qml.Object
}

var engine *qml.Engine
var win *qml.Window

func main() {
	if err := qml.SailfishRun("harbour-whisperfish", "", VERSION, run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runBackend() {
	refreshContacts()
}

func run() error {
	whisperfish := Whisperfish{}
	engine = qml.SailfishNewEngine()
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
