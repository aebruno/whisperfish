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

package tools

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/whisperfish/client"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/sailfish"
	"golang.org/x/crypto/ssh/terminal"
)

func ConvertDataStore() {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get password from terminal")
	}
	password := strings.TrimSpace(string(bytePassword))
	fmt.Println()

	log.Infof("Converting data store: %s", password)
	app := sailfish.SailfishApp_Application(len(os.Args), os.Args)
	app.SetOrganizationName("")
	app.SetApplicationName("harbour-whisperfish")

	var configPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__ConfigLocation)
	var dataPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__DataLocation)

	var backend = client.NewBackend(nil)
	err = backend.Setup(configPath, dataPath, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to setup backend")
	}
}
