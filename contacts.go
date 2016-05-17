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

package main

import (
	"github.com/aebruno/qml"
	"github.com/janimo/textsecure"
	"github.com/ttacon/libphonenumber"
)

type Contacts struct {
	contacts []textsecure.Contact
	Len      int
}

// Get contact by index i
func (c *Contacts) Contact(i int) textsecure.Contact {
	if i == -1 {
		return textsecure.Contact{}
	}
	return c.contacts[i]
}

// Get name of contact with number tel
func (c *Contacts) Name(tel string) string {
	for _, r := range c.contacts {
		if r.Tel == tel {
			return r.Name
		}
	}

	// name not found. just return number
	return tel
}

// Find contact by tel
func (c *Contacts) Find(tel, countryCode string) textsecure.Contact {
	num, err := libphonenumber.Parse(tel, countryCode)
	if err != nil {
		return textsecure.Contact{}
	}

	n := libphonenumber.Format(num, libphonenumber.E164)
	for i, r := range c.contacts {
		if r.Tel == n {
			return c.contacts[i]
		}
	}

	return textsecure.Contact{}
}

// Refresh list of local contacts that are registered with Signal
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
