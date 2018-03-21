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

	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/settings"
	"github.com/aebruno/whisperfish/store"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/ttacon/libphonenumber"
)

//go:generate qtmoc
type ContactModel struct {
	core.QAbstractListModel

	contacts           []*store.Contact
	registeredContacts []textsecure.Contact
	settings           *settings.Settings

	_ map[int]*core.QByteArray     `property:"roles"`
	_ func()                       `constructor:"init"`
	_ func()                       `signal:"refreshComplete"`
	_ func(tel string) string      `slot:"format"`
	_ func(tel string) bool        `slot:"registered"`
	_ func(tel string) string      `slot:"identity"`
	_ func(tel string) string      `slot:"name"`
	_ func() int                   `slot:"total"`
	_ func()                       `slot:"refresh"`
	_ func(row int) *core.QVariant `slot:"get"`
}

// Setup connections
func (model *ContactModel) init() {
	model.settings = settings.NewSettings(nil)
	model.contacts = make([]*store.Contact, 0)
	model.registeredContacts = make([]textsecure.Contact, 0)

	model.SetRoles(map[int]*core.QByteArray{
		RoleName: core.NewQByteArray2("name", len("name")),
		RoleTel:  core.NewQByteArray2("tel", len("tel")),
	})

	// Slot connections
	model.ConnectRoleNames(model.roleNames)
	model.ConnectData(model.data)
	model.ConnectColumnCount(model.columnCount)
	model.ConnectRowCount(model.rowCount)
	model.ConnectIdentity(model.identity)
	model.ConnectFormat(model.format)
	model.ConnectRegistered(model.registered)
	model.ConnectName(model.name)
	model.ConnectTotal(model.total)
	model.ConnectRefresh(model.refresh)
	model.ConnectGet(model.get)
}

// Returns the data stored under the given role for the item referred to by the
// index.
func (model *ContactModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() < 0 || index.Row() > len(model.contacts) {
		return core.NewQVariant()
	}

	contact := model.contacts[index.Row()]
	switch role {
	case RoleName:
		return core.NewQVariant14(contact.Name)
	case RoleTel:
		return core.NewQVariant14(contact.Tel)
	default:
		return core.NewQVariant()
	}
}

// Returns the number of items in the model.
func (model *ContactModel) rowCount(parent *core.QModelIndex) int {
	return len(model.contacts)
}

// Return the number of columns. This will always be 1
func (model *ContactModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

// Return the roles for the model
func (model *ContactModel) roleNames() map[int]*core.QByteArray {
	return model.Roles()
}

// Get the item at index
func (model *ContactModel) get(row int) *core.QVariant {
	rec := make(map[string]*core.QVariant)

	if row < 0 || row > len(model.contacts) {
		return core.NewQVariant()
	}

	contact := model.contacts[row]
	rec["name"] = core.NewQVariant14(contact.Name)
	rec["tel"] = core.NewQVariant14(contact.Tel)

	return core.NewQVariant25(rec)
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
