package main

import (
	"github.com/janimo/textsecure"
	"gopkg.in/qml.v1"
)

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

var contactsModel *Contacts = &Contacts{}

func refreshContacts() error {
	var err error
	contactsModel.contacts, err = getAddressBookContacts()
	if err != nil {
		return err
	}
	contactsModel.Len = len(contactsModel.contacts)
	qml.Changed(contactsModel, &contactsModel.Len)

	return nil
}

func initModels() {
	engine.Context().SetVar("contactsModel", contactsModel)
}
