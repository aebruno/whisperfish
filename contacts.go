package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/janimo/textsecure"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ttacon/libphonenumber"
	"gopkg.in/qml.v1"
)

const QTCONTACTS_PATH = "/home/nemo/.local/share/system/Contacts/qtcontacts-sqlite/contacts.db"

type Contacts struct {
	contacts []textsecure.Contact
	Len      int
}

func (c *Contacts) Contact(i int) textsecure.Contact {
	if i == -1 {
		return textsecure.Contact{}
	}
	return c.contacts[i]
}

func (c *Contacts) Refresh() error {
	signalContacts, err := textsecure.GetRegisteredContacts()
	if err != nil {
		return err
	}

	c.contacts = signalContacts
	c.Len = len(c.contacts)
	qml.Changed(c, &c.Len)

	return nil
}

func (c *Contacts) Init() error {
	var err error
	c.contacts, err = getSailfishContacts()
	if err != nil {
		return err
	}

	c.Len = len(c.contacts)
	qml.Changed(c, &c.Len)

	return nil
}

func getSailfishContacts() ([]textsecure.Contact, error) {
	db, err := sqlx.Open("sqlite3", QTCONTACTS_PATH)
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
		n = strings.TrimPrefix(n, "+")
		num, err := libphonenumber.Parse(fmt.Sprintf("+%s", n), "")
		if err == nil {
			contacts[i].Tel = libphonenumber.Format(num, libphonenumber.E164)
		}
	}

	return contacts, nil
}
