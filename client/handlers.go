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

package client

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/store"
	"github.com/ttacon/libphonenumber"
)

// Prompt the user for storage password and create encrypted data store if
// needed
func (b *Backend) getStoragePassword() string {
	pass := b.prompt.GetStoragePassword()

	if b.settings.GetBool("encrypt_database") {
		err := b.newStorage(pass)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to open encrypted database")
		}
	}

	return pass
}

// Message handler
func (b *Backend) messageHandler(msg *textsecure.Message) {
	b.processMessage(msg, false, 0)
}

// Receipt handler
func (b *Backend) receiptHandler(source string, devID uint32, ts uint64) {
	log.WithFields(log.Fields{
		"source":    source,
		"timestamp": ts,
		"devID":     devID,
	}).Debug("Receipt handler")

	var err error
	sessionID := int64(0)
	messageID := int64(0)
	tries := 0

	for {
		sessionID, messageID, err = b.ds.MarkMessageReceived(source, ts)
		if err != nil {
			tries++
			if tries > 3 {
				log.WithFields(log.Fields{
					"error":     err,
					"source":    source,
					"timestamp": ts,
				}).Error("Failed to mark message received")
				return
			}
			log.Debug("receiptHandler can't find message. Trying again later")
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	err = b.ds.MarkSessionReceived(sessionID, ts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    sessionID,
		}).Error("Failed to mark session received")
	}

	b.sessionModel.MarkReceived(sessionID)
	b.messageModel.MarkReceived(messageID)
}

// Registration handler
func (b *Backend) registrationDone() {
	textsecure.WriteConfig(b.configFile, b.config)

	num, err := libphonenumber.Parse(b.config.Tel, "")
	if err == nil {
		b.settings.SetString("country_code", libphonenumber.GetRegionCodeForNumber(num))
	}

	log.Info("Registered")
	b.RegistrationSuccess()
	if _, err := os.Stat(filepath.Join(b.config.StorageDir, "identity", "identity_key")); err == nil {
		b.SetRegistered(true)
	}
}

func (b *Backend) syncReadHandler(source string, ts uint64) {
	log.Info("TODO: Processing sync read message")
}

func (b *Backend) syncSentHandler(msg *textsecure.Message, ts uint64) {
	log.Info("TODO: Processing sync sent message")
}

func (b *Backend) getLocalContacts() ([]textsecure.Contact, error) {
	return store.SailfishContacts(b.settings.GetString("country_code"))
}
