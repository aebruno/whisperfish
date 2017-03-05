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

package model

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/settings"
	"github.com/aebruno/whisperfish/store"
	"github.com/therecipe/qt/core"
	"github.com/ttacon/libphonenumber"
)

//go:generate qtmoc
type ContactModel struct {
	core.QObject

	contacts           []*store.Contact
	registeredContacts []textsecure.Contact
	settings           *settings.Settings

	_ func()                  `constructor:"init"`
	_ func()                  `signal:"refreshComplete"`
	_ func(tel string) string `slot:"format"`
	_ func(tel string) bool   `slot:"registered"`
	_ func(tel string) string `slot:"identity"`
	_ func(tel string) string `slot:"name"`
	_ func() int              `slot:"total"`
	_ func()                  `slot:"refresh"`
}

// Setup connections
func (model *ContactModel) init() {
	model.settings = settings.NewSettings(nil)
	model.contacts = make([]*store.Contact, 0)
	model.registeredContacts = make([]textsecure.Contact, 0)

	// Slot connections
	model.ConnectIdentity(model.identity)
	model.ConnectFormat(model.format)
	model.ConnectRegistered(model.registered)
	model.ConnectName(model.name)
	model.ConnectTotal(model.total)
	model.ConnectRefresh(model.refresh)
}

// Format contact number
func (model *ContactModel) format(tel string) string {
	num, err := libphonenumber.Parse(tel, model.settings.GetString("country_code"))
	if err != nil {
		return ""
	}

	n := libphonenumber.Format(num, libphonenumber.E164)
	return n
}

// Returns the name of the contact
func (model *ContactModel) name(tel string) string {
	for _, r := range model.contacts {
		if r.Tel == tel {
			return r.Name
		}
	}

	// name not found. just return number
	return tel
}

// Returns the total number of contacts registered with signal
func (model *ContactModel) total() int {
	return len(model.registeredContacts)
}

// Refresh contacts
func (model *ContactModel) refresh() {
	var err error
	model.contacts, err = store.SailfishContacts(model.settings.GetString("country_code"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch local contacts")
		model.contacts = make([]*store.Contact, 0)
	}

	if model.settings.GetBool("share_contacts") {
		model.registeredContacts, err = textsecure.GetRegisteredContacts()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to fetch signal contacts")
			model.registeredContacts = make([]textsecure.Contact, 0)
		}
	}
}

// Check if contact is registered with signal
func (model *ContactModel) registered(tel string) bool {
	num, err := libphonenumber.Parse(tel, model.settings.GetString("country_code"))
	if err != nil {
		return false
	}

	n := libphonenumber.Format(num, libphonenumber.E164)

	for _, r := range model.registeredContacts {
		if r.Tel == n {
			return true
		}
	}

	return false
}

// Returns identity of contact
func (model *ContactModel) identity(source string) string {
	id, err := textsecure.ContactIdentityKey(source)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch contact identity")
		return ""
	}

	return fmt.Sprintf("% 0X", id)
}
