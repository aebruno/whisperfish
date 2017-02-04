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

package worker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/store"
	"github.com/janimo/textsecure/axolotl"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type SendWorker struct {
	core.QObject

	config *textsecure.Config

	_ func()                               `constructor:"init"`
	_ func(sid int64)                      `signal:"sendMessage"`
	_ func(sid, mid int64, message string) `signal:"messageSent"`
	_ func(source string)                  `signal:"promptResetPeerIdentity"`
	_ func(confirm string)                 `signal:"resetPeerIdentity"`
}

func (s *SendWorker) init() {
	s.ConnectSendMessage(func(sid int64) {
		go s.sendMessage(sid)
	})
}

func (s *SendWorker) SetConfig(config *textsecure.Config) {
	s.config = config
}

// Send message to Signal server
func (s *SendWorker) send(m *store.Message) error {
	var att io.Reader
	var err error
	var ts uint64

	sess, err := store.DS.FetchSession(m.SID)
	if err != nil {
		return err
	}

	if m.Attachment != "" {
		att, err = os.Open(m.Attachment)
		if err != nil {
			return err
		}
	}

	if m.Flags == textsecure.EndSessionFlag {
		ts, err = textsecure.EndSession(sess.Source, "Reset Secure Session")
	} else if att == nil {
		if sess.IsGroup {
			ts, err = textsecure.SendGroupMessage(sess.Source, m.Message)
		} else {
			ts, err = textsecure.SendMessage(sess.Source, m.Message)
		}
	} else {
		if sess.IsGroup {
			ts, err = textsecure.SendGroupAttachment(sess.Source, m.Message, att)
		} else {
			ts, err = textsecure.SendAttachment(sess.Source, m.Message, att)
		}
	}

	if nerr, ok := err.(axolotl.NotTrustedError); ok {
		remoteIdentityPath := filepath.Join(s.config.StorageDir, "identity", fmt.Sprintf("remote_%s", nerr.ID))
		log.WithFields(log.Fields{
			"error":          err,
			"source":         nerr.ID,
			"remoteIdentity": remoteIdentityPath,
		}).Error("Peer identity not trusted")

		confirm := s.getConfirmResetPeerIdentity(nerr.ID)
		if confirm == "yes" {
			err = os.Remove(remoteIdentityPath)
			if err != nil {
				return err
			}

			// retry
			go s.sendMessage(m.ID)

			return fmt.Errorf("Reset peer identity")
		}

		err := store.DS.DequeueSent(m.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"id":    m.ID,
			}).Error("Failed to remove message from mailq")
		}
		return fmt.Errorf("Peer identity not trusted. Abort sending message.")
	}

	if err != nil {
		return err
	}

	err = store.DS.MarkSessionSent(sess.ID, m.Message, ts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    sess.ID,
		}).Error("Failed to mark session sent")
		return err
	}

	err = store.DS.MarkMessageSent(m.ID, ts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    m.ID,
		}).Error("Failed to mark message sent")
		return err
	}

	s.MessageSent(sess.ID, m.ID, m.Message)

	return nil
}

// Fetch message from queue and send
func (s *SendWorker) sendMessage(id int64) {
	message, err := store.DS.FetchQueuedMessage(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to fetch message from queue")
	}

	err = s.send(message)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to send message")
		return
	}

	// Remove from sentq
	err = store.DS.DequeueSent(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to remove message from queue")
	}
}

// Prompt the user to confirm reset peer identity
func (s *SendWorker) getConfirmResetPeerIdentity(source string) string {
	log.Infof("Prompting to reset peer identity: %s", source)
	ch := make(chan string)
	s.ConnectResetPeerIdentity(func(confirm string) {
		ch <- confirm
	})

	s.PromptResetPeerIdentity(source)
	confirm := <-ch
	return confirm
}
