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

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/whisperfish/client"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/sailfish"
)

// Run the main QT gui thread
func Run(version string) {
	log.Infof("Starting Whisperfish version %s", version)
	app := sailfish.SailfishApp_Application(len(os.Args), os.Args)
	app.SetApplicationVersion(version)
	app.SetOrganizationName("")
	app.SetApplicationName("harbour-whisperfish")

	var view = sailfish.SailfishApp_CreateView()

	var backend = client.NewBackend(nil)
	backend.ConnectCopyToClipboard(func(text string) {
		log.Info("Copy to clipboard")
		if len(text) > 0 {
			app.Clipboard().Clear(gui.QClipboard__Clipboard)
			app.Clipboard().SetText(text, gui.QClipboard__Clipboard)
		}
	})

	var configPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__ConfigLocation)
	var dataPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__DataLocation)

	err := backend.Setup(configPath, dataPath, view)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to setup backend")
	}

	view.SetSource(sailfish.SailfishApp_PathTo("qml/harbour-whisperfish.qml"))
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.Show()

	go backend.Run()

	gui.QGuiApplication_Exec()
}
