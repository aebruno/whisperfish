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
	"mime"
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

func AddAttachmentExtensions() {
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
	fmt.Println("WARNING: This operation will modify your attachment files. Ensure you have a backup!")
	fmt.Println("**************************")
	fmt.Print("Do you want to proceed? Type yes or no: ")
	scanner.Scan()
	ans := scanner.Text()
	ans = strings.TrimSpace(ans)
	if ans != "yes" {
		os.Exit(0)
	}

	password := ""
	if settings.GetBool("encrypt_database") {
		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Failed to get password from terminal")
		}
		password = strings.TrimSpace(string(bytePassword))
		fmt.Println()

		if len(password) < 6 {
			log.Fatalf("Password must be >6 characters long")
		}
	}

	log.Infof("Adding extensions to attachments")

	err = processAttachments(settings, dataPath, password)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to process attachments")
	}

	log.Info("Attachments processed succesfully")
}

func processAttachments(settings *settings.Settings, dataPath, password string) error {
	ds, err := NewStorage(dataPath, password)
	if err != nil {
		return err
	}

	type record struct {
		MessageID  int64  `db:"id"`
		Attachment string `db:"attachment"`
		MimeType   string `db:"mime_type"`
	}

	records := []*record{}
	err = ds.DBX().Select(&records, `
        select id, attachment, mime_type from message where has_attachment = 1 and outgoing = 0
    `)
	if err != nil {
		return err
	}

	total := 0
	cnt := 0
	for _, rec := range records {
		name := strings.TrimSuffix(rec.Attachment, filepath.Ext(rec.Attachment))
		ext, _ := mime.ExtensionsByType(rec.MimeType)
		if ext == nil {
			ext = []string{""}
		}

		fname := fmt.Sprintf("%s%s", name, ext[0])

		log.WithFields(log.Fields{
			"id":        rec.MessageID,
			"origPath":  rec.Attachment,
			"newPath":   fname,
			"extension": ext[0],
		}).Info("Renaming attachment file")

		err := os.Rename(rec.Attachment, fname)
		if err != nil {
			log.WithFields(log.Fields{
				"id":        rec.MessageID,
				"origPath":  rec.Attachment,
				"newPath":   fname,
				"extension": ext[0],
			}).Error("Failed to rename attachment file")
			total++
			continue
		}
		_, err = ds.DBX().Exec(`update message set attachment = ? where id = ?`, fname, rec.MessageID)
		if err != nil {
			log.WithFields(log.Fields{
				"id":        rec.MessageID,
				"origPath":  rec.Attachment,
				"newPath":   fname,
				"extension": ext[0],
			}).Error("Failed to update attachment path in database")
			total++
			continue
		}
		cnt++
		total++
	}

	log.Infof("Processed %d of %d attachments successfully", cnt, total)

	return nil
}
