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
	"bufio"
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
	app := sailfish.SailfishApp_Application(len(os.Args), os.Args)
	app.SetOrganizationName("")
	app.SetApplicationName("harbour-whisperfish")

	var configPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__ConfigLocation)
	var dataPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__DataLocation)

	var backend = client.NewBackend(nil)
	err := backend.Setup(configPath, dataPath, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to setup backend")
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("***** DANGER ZONE ********")
	fmt.Println("WARNING: This operation could lock you out of using Whisperfish.")
	fmt.Println()
	if backend.HasEncryptedDatabase() {
		fmt.Println("Your database is currently encrypted. This operation will")
		fmt.Println("decrypt the database. This is NOT recommended and should only")
		fmt.Println("be used for development purposes.")
	} else {
		fmt.Println("Your database is currently decrypted. This operation this operation")
		fmt.Println("will encrypt your database.")
	}
	fmt.Println("**************************")
	fmt.Print("Do you want to proceed? Type yes or no: ")
	scanner.Scan()
	ans := scanner.Text()
	ans = strings.TrimSpace(ans)
	if ans != "yes" {
		os.Exit(0)
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get password from terminal")
	}
	password := strings.TrimSpace(string(bytePassword))
	fmt.Println()

	if len(password) < 6 {
		log.Fatalf("Password must be >6 characters long")
	}

	log.Infof("Converting data store")

	err = backend.ConvertDataStore(password)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to convert datastore")
	}

	log.Info("Data store converted successfully")
}
