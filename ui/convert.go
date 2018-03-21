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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/aebruno/whisperfish/settings"
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

	configDir, storageDir := InitDirs(configPath, dataPath)

	var settings = settings.NewSettings(nil)
	err := settings.Setup(configDir, storageDir)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to initialize settings")
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("***** DANGER ZONE ********")
	fmt.Println("WARNING: This operation could lock you out of using Whisperfish.")
	fmt.Println()
	if settings.GetBool("encrypt_database") {
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

	err = convert(settings, dataPath, password)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to convert datastore")
	}

	log.Info("Data store converted successfully")
}

func convert(settings *settings.Settings, dataPath, password string) error {
	if password == "" {
		return fmt.Errorf("No password given")
	}

	dbDir := filepath.Join(dataPath, "db")
	tmp := filepath.Join(dbDir, "tmp.db")
	encrypt := !settings.GetBool("encrypt_database")

	if encrypt {
		log.Info("Encrypting database..")

		ds, err := NewStorage(dataPath, "")
		if err != nil {
			return err
		}

		err = ds.Encrypt(tmp, password)
		if err != nil {
			return err
		}
	} else {
		log.Info("Decrypting database..")

		ds, err := NewStorage(dataPath, password)
		if err != nil {
			return err
		}

		err = ds.Decrypt(tmp)
		if err != nil {
			return err
		}
	}

	dbFile := filepath.Join(dbDir, WhisperDB)

	err := os.Rename(tmp, dbFile)
	if err != nil {
		return err
	}

	settings.SetBool("encrypt_database", encrypt)
	settings.Sync()
	return nil
}
