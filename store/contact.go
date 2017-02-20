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
	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/ttacon/libphonenumber"
)

const (
	// Path to Sailfish read-only contacts database
	QtcontactsPath = "/home/nemo/.local/share/system/Contacts/qtcontacts-sqlite/contacts.db"
)

type Contact struct {
	contacts []textsecure.Contact
}

func NewContact() *Contact {
	c := &Contact{}
	return c
}

func (c *Contact) Len() int {
	return len(c.contacts)
}

func (c *Contact) RegisteredContacts() int {
	contacts, err := textsecure.GetRegisteredContacts()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch signal contacts")
		return 0
	}

	count := 0
	for _, l := range c.contacts {
		for _, r := range contacts {
			if l.Tel == r.Tel {
				count++
				break
			}
		}
	}

	return count
}

func (c *Contact) Refresh(country string) {
	contacts, err := SailfishContacts(country)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to refresh contacts")
		c.contacts = make([]textsecure.Contact, 0)
	}

	c.contacts = contacts
}

// Get name of contact with number tel
func (c *Contact) FindName(tel string) string {
	for _, r := range c.contacts {
		if r.Tel == tel {
			return r.Name
		}
	}

	// name not found. just return number
	return tel
}

// Format contact tel
func (c *Contact) Format(tel, countryCode string) string {
	num, err := libphonenumber.Parse(tel, countryCode)
	if err != nil {
		return ""
	}

	n := libphonenumber.Format(num, libphonenumber.E164)
	return n
}

// Check if contact is registered with Signal
func (c *Contact) Exists(tel, countryCode string) bool {
	num, err := libphonenumber.Parse(tel, countryCode)
	if err != nil {
		return false
	}

	n := libphonenumber.Format(num, libphonenumber.E164)

	contacts, err := textsecure.GetRegisteredContacts()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch signal contacts")
		return false
	}

	for _, r := range contacts {
		if r.Tel == n {
			return true
		}
	}

	return false
}

// Get local sailfish contacts
func SailfishContacts(country string) ([]textsecure.Contact, error) {
	db, err := sqlx.Open("sqlite3", QtcontactsPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to open contacts database")
		return nil, err
	}

	contacts := []textsecure.Contact{}
	err = db.Select(&contacts, `
	select
	   c.displayLabel as name,
	   p.phoneNumber as tel
	from Contacts as c
	join PhoneNumbers p
	on c.contactId = p.contactId`)
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
