package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/janimo/textsecure"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func getAddressBookContacts() ([]textsecure.Contact, error) {
	db, err := sqlx.Open("sqlite3", "/home/nemo/.local/share/system/Contacts/qtcontacts-sqlite/contacts.db")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to open contacts database")
		return nil, err
	}

	contacts := []textsecure.Contact{}
	err = db.Select(&contacts, "select c.displayLabel as name,p.phoneNumber as tel from Contacts as c join PhoneNumbers p on c.contactId = p.contactId")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to query contacts database")
		return nil, err
	}

	return contacts, nil
}
