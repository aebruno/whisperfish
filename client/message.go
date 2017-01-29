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

package client

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/store"
	"github.com/janimo/textsecure/axolotl"
)

// Process incoming message
func (b *Backend) processMessage(msg *textsecure.Message, isSyncSent bool, ts uint64) {
	log.Infof("Received message from: %s", msg.Source())

	message := &store.Message{
		Source:  msg.Source(),
		Message: msg.Message(),
		Flags:   msg.Flags(),
	}

	if isSyncSent {
		message.Outgoing = true
		message.Sent = true
		if ts > 0 {
			message.Timestamp = ts
		}
	} else {
		message.Timestamp = msg.Timestamp()
	}

	if len(msg.Attachments()) > 0 {
		if b.settings.GetBool("save_attachments") && !b.settings.GetBool("incognito") {
			err := message.SaveAttachment(b.settings.GetString("attachment_dir"), msg.Attachments()[0])
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed to save attachment")
			}
		} else {
			message.HasAttachment = true
			message.MimeType = msg.Attachments()[0].MimeType
		}
	}

	err := b.updateModel(message, msg.Group(), !message.Sent)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add/update session in database")
		return
	}

	// Don't send notification if disabled
	if !b.settings.GetBool("enable_notify") {
		return
	}

	b.NotifyMessage(message.SID, b.contacts.FindName(msg.Source()), msg.Message())
}

// Send message to Signal server
func (b *Backend) sendSignalMessage(m *store.Message) error {
	var att io.Reader
	var err error
	var ts uint64

	s, err := b.ds.FetchSession(m.SID)
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
		ts, err = textsecure.EndSession(s.Source, "Reset Secure Session")
	} else if att == nil {
		if s.IsGroup {
			ts, err = textsecure.SendGroupMessage(s.Source, m.Message)
		} else {
			ts, err = textsecure.SendMessage(s.Source, m.Message)
		}
	} else {
		if s.IsGroup {
			ts, err = textsecure.SendGroupAttachment(s.Source, m.Message, att)
		} else {
			ts, err = textsecure.SendAttachment(s.Source, m.Message, att)
		}
	}

	if nerr, ok := err.(axolotl.NotTrustedError); ok {
		remoteIdentityPath := filepath.Join(b.config.StorageDir, "identity", fmt.Sprintf("remote_%s", nerr.ID))
		log.WithFields(log.Fields{
			"error":          err,
			"source":         nerr.ID,
			"remoteIdentity": remoteIdentityPath,
		}).Error("Peer identity not trusted")

		confirm := b.prompt.GetConfirmResetPeerIdentity(nerr.ID)
		if confirm == "yes" {
			err = os.Remove(remoteIdentityPath)
			if err != nil {
				return err
			}

			return fmt.Errorf("Reset peer identity")
		}

		err := b.ds.DequeueSent(m.ID)
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

	err = b.ds.MarkSessionSent(s.ID, m.Message, ts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    s.ID,
		}).Error("Failed to mark session sent")
		return err
	}

	err = b.ds.MarkMessageSent(m.ID, ts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    m.ID,
		}).Error("Failed to mark message sent")
		return err
	}

	b.sessionModel.MarkSent(s.ID)
	b.messageModel.MarkSent(m.ID)

	return nil
}

// Worker thread that checks message queue and sends message
func (b *Backend) sendMessageWorker() {
	for {
		time.Sleep(3 * time.Second)

		if !b.IsConnected() {
			continue
		}

		messages, err := b.ds.FetchSentq()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to fetch mailq")
		}

		for _, m := range messages {
			log.Debugf("Sending message: %d", m.ID)

			err = b.sendSignalMessage(m)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    m.ID,
				}).Error("Failed to send message")
				continue
			}

			// Remove from sentq
			err := b.ds.DequeueSent(m.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    m.ID,
				}).Error("Failed to remove message from mailq")
			}

			// Throttle
			time.Sleep(1 * time.Second)
		}
	}
}

// Reset secure session
func (b *Backend) endSession(source string) {
	message := &store.Message{
		Source:    source,
		Message:   "[Whisperfish] Reset secure session",
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Outgoing:  true,
		Flags:     textsecure.EndSessionFlag,
	}

	err := b.updateModel(message, nil, false)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add EndSession message to database")
		return
	}

	err = b.ds.QueueSent(message)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add EndSession message to queue")
	}
}

// Queue outgoing message for delivery to Signal
func (b *Backend) queueMessage(to, msg, attachment string, group *textsecure.Group) (int64, error) {
	message := &store.Message{
		Source:    to,
		Message:   msg,
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Outgoing:  true,
	}

	if len(attachment) > 0 {
		att, err := os.Open(attachment)
		if err != nil {
			return 0, err
		}
		defer att.Close()
		//XXX Sucks we have to do this twice
		message.MimeType, _ = textsecure.MIMETypeFromReader(att)
		message.Attachment = attachment
		message.HasAttachment = true
	}

	err := b.updateModel(message, group, false)
	if err != nil {
		return 0, err
	}

	err = b.ds.QueueSent(message)
	if err != nil {
		return 0, err
	}

	return message.SID, nil
}

// Perpare outgoing message for delivery to Signal
func (b *Backend) sendMessage(source, message, groupName, attachment string) int64 {
	var err error
	var sid int64

	m := strings.Split(source, ",")
	if len(m) > 1 {
		group, err := textsecure.NewGroup(groupName, m)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"group_name": groupName,
			}).Error("Failed to create new group")
			return 0
		}

		sid, err = b.queueMessage(group.Hexid, message, attachment, group)
	} else {
		sid, err = b.queueMessage(source, message, attachment, nil)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   b.messageModel.SessionId(),
		}).Error("Failed to send message")
	}

	return sid
}

func (b *Backend) updateModel(message *store.Message, group *textsecure.Group, unread bool) error {
	err := b.sessionModel.ProcessMessage(message, group, unread)
	if err != nil {
		return err
	}

	err = b.messageModel.SaveMessage(message)
	if err != nil {
		return err
	}

	return nil
}
