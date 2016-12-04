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
	"io/ioutil"

	"github.com/therecipe/qt/core"
	"gopkg.in/yaml.v2"
)

//go:generate qtmoc
type Settings struct {
	core.QObject

	_ func(key string, val string) `slot:"stringSet"`
	_ func(key string, val bool)   `slot:"boolSet"`
	_ func(key string) string      `slot:"stringValue"`
	_ func(key string) bool        `slot:"boolValue"`
	_ func()                       `slot:"defaults"`

	_ func() `constructor:"init"`
}

// Initialize settings and setup slot connections
func (s *Settings) init() {
	// Connect slots
	s.ConnectStringSet(func(key string, val string) {
		s.SetString(key, val)
	})
	s.ConnectBoolSet(func(key string, val bool) {
		s.SetBool(key, val)
	})
	s.ConnectStringValue(func(key string) string {
		return s.GetString(key)
	})
	s.ConnectBoolValue(func(key string) bool {
		return s.GetBool(key)
	})
	s.ConnectDefaults(func() {
		s.SetDefaults()
	})
}

// Set string value
func (s *Settings) SetString(key string, val string) {
	var settings = core.NewQSettings6(nil)
	settings.SetValue(key, core.NewQVariant18(val))
}

// Set bool value
func (s *Settings) SetBool(key string, val bool) {
	var settings = core.NewQSettings6(nil)
	settings.SetValue(key, core.NewQVariant12(val))
}

// Get string value
func (s *Settings) GetString(key string) string {
	var settings = core.NewQSettings6(nil)
	return settings.Value(key, core.NewQVariant18("")).ToString()
}

// Get bool value
func (s *Settings) GetBool(key string) bool {
	var settings = core.NewQSettings6(nil)
	return settings.Value(key, core.NewQVariant12(false)).ToBool()
}

// Set default values
func (s *Settings) SetDefaults() {
	var settings = core.NewQSettings6(nil)
	settings.SetValue("incognito", core.NewQVariant12(false))
	settings.SetValue("enable_notify", core.NewQVariant12(true))
	settings.SetValue("show_notify_message", core.NewQVariant12(false))
	settings.SetValue("encrypt_database", core.NewQVariant12(true))
	settings.SetValue("save_attachments", core.NewQVariant12(true))
	settings.SetValue("country_code", core.NewQVariant18(""))
	settings.Sync()
}

func (s *Settings) MigrateSettings(file string) error {

	type deprecatedSettings struct {
		Incognito         bool   `yaml:"incognito"`
		EnableNotify      bool   `yaml:"enable_notify"`
		ShowNotifyMessage bool   `yaml:"show_notify_message"`
		EncryptDatabase   bool   `yaml:"encrypt_database"`
		SaveAttachments   bool   `yaml:"save_attachments"`
		CountryCode       string `yaml:"country_code"`
	}

	ds := &deprecatedSettings{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, ds)
	if err != nil {
		return err
	}

	var settings = core.NewQSettings6(nil)
	settings.SetValue("incognito", core.NewQVariant12(ds.Incognito))
	settings.SetValue("enable_notify", core.NewQVariant12(ds.EnableNotify))
	settings.SetValue("show_notify_message", core.NewQVariant12(ds.ShowNotifyMessage))
	settings.SetValue("encrypt_database", core.NewQVariant12(ds.EncryptDatabase))
	settings.SetValue("save_attachments", core.NewQVariant12(ds.SaveAttachments))
	settings.SetValue("country_code", core.NewQVariant18(ds.CountryCode))

	settings.Sync()

	return nil
}

func (s *Settings) Sync() {
	var settings = core.NewQSettings6(nil)
	settings.Sync()
}
