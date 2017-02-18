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
)

//go:generate qtmoc
type Contact struct {
	core.QObject

	settings     *settings.Settings
	contactStore *store.Contact

	_ func()                  `constructor:"init"`
	_ func()                  `signal:"refreshComplete"`
	_ func(tel string) string `slot:"format"`
	_ func(tel string) bool   `slot:"exists"`
	_ func(tel string) string `slot:"identity"`
	_ func(tel string) string `slot:"name"`
	_ func() int              `slot:"total"`
	_ func()                  `slot:"refresh"`
}

// Setup connections
func (c *Contact) init() {
	c.settings = settings.NewSettings(nil)
	c.contactStore = store.NewContact()

	// Slot connections
	c.ConnectIdentity(func(source string) string {
		return c.identity(source)
	})
	c.ConnectFormat(func(tel string) string {
		return c.contactStore.Format(tel, c.settings.GetString("country_code"))
	})
	c.ConnectExists(func(tel string) bool {
		return c.contactStore.Exists(tel, c.settings.GetString("country_code"))
	})
	c.ConnectName(func(tel string) string {
		return c.contactStore.FindName(tel)
	})
	c.ConnectTotal(func() int {
		return c.contactStore.Len()
	})
	c.ConnectRefresh(func() {
		c.contactStore.Refresh(c.settings.GetString("country_code"))
		c.RefreshComplete()
	})
}

// Returns identity of contact
func (c *Contact) identity(source string) string {
	id, err := textsecure.ContactIdentityKey(source)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch contact identity")
		return ""
	}

	return fmt.Sprintf("% 0X", id)
}
