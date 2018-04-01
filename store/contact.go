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

package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher"
	log "github.com/sirupsen/logrus"
	"github.com/ttacon/libphonenumber"
)

const (
	// Path to Sailfish read-only contacts database
	QtcontactsPath = "/home/nemo/.local/share/system/Contacts/qtcontacts-sqlite/contacts.db"
)

type Contact struct {
	Name string `db:"name"`
	Tel  string `db:"tel"`
}

// Get local sailfish contacts
func SailfishContacts(country string) ([]*Contact, error) {
	db, err := sqlx.Open("sqlite3", QtcontactsPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to open contacts database")
		return nil, err
	}

	contacts := []*Contact{}
	err = db.Select(&contacts, `
	select
	   c.displayLabel as name,
	   p.phoneNumber as tel
	from Contacts as c
	join PhoneNumbers p
	on c.contactId = p.contactId
    order by name`)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to query contacts database")
		return nil, err
	}

	// Reformat numbers in E.164 format
	for i := range contacts {
		n := contacts[i].Tel
		num, err := libphonenumber.Parse(n, country)
		if err == nil {
			contacts[i].Tel = libphonenumber.Format(num, libphonenumber.E164)
		}
	}

	return contacts, nil
}
