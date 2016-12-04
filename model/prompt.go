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
	"github.com/therecipe/qt/core"
	"github.com/ttacon/libphonenumber"
)

//go:generate qtmoc
type Prompt struct {
	core.QObject

	_ func()              `signal:"promptPhoneNumber"`
	_ func()              `signal:"promptVerificationCode"`
	_ func()              `signal:"promptPassword"`
	_ func(source string) `signal:"promptResetPeerIdentity"`

	_ func(number string)  `slot:"phoneNumber"`
	_ func(code string)    `slot:"verificationCode"`
	_ func(passwd string)  `slot:"password"`
	_ func(confirm string) `slot:"resetPeerIdentity"`
}

// Prompt the user to enter telephone number for Registration
func (p *Prompt) GetPhoneNumber() string {
	log.Info("Prompting for phone number")
	ch := make(chan string)
	p.ConnectPhoneNumber(func(number string) {
		ch <- number
	})

	p.PromptPhoneNumber()
	n := <-ch
	num, err := libphonenumber.Parse(fmt.Sprintf("+%s", n), "")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to parse phone number")
	}

	tel := libphonenumber.Format(num, libphonenumber.E164)
	log.Infof("Using phone number: %s", tel)
	return tel
}

// Prompt the user to enter the verification code
func (p *Prompt) GetVerificationCode() string {
	log.Info("Prompting for verification code")
	ch := make(chan string)
	p.ConnectVerificationCode(func(code string) {
		ch <- code
	})

	p.PromptVerificationCode()
	code := <-ch
	log.Infof("Code: %s", code)
	return code
}

// Prompt the user for storage password
func (p *Prompt) GetStoragePassword() string {
	log.Info("Prompting for storage password")
	ch := make(chan string)
	p.ConnectPassword(func(passwd string) {
		ch <- passwd
	})

	p.PromptPassword()
	pass := <-ch

	return pass
}

// Prompt the user to confirm reset peer identity
func (p *Prompt) GetConfirmResetPeerIdentity(source string) string {
	log.Infof("Prompting to reset peer identity: %s", source)
	ch := make(chan string)
	p.ConnectResetPeerIdentity(func(confirm string) {
		ch <- confirm
	})

	p.PromptResetPeerIdentity(source)
	confirm := <-ch
	return confirm
}
