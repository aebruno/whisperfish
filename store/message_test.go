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
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aebruno/textsecure"
)

func TestMessage(t *testing.T) {
	ds, err := newTestDataStore()
	if err != nil {
		t.Error(err)
	}

	now := uint64(time.Now().UnixNano() / 1000000)

	tel := "+1771111006"
	text := "Hello World"
	sess := &Session{Source: tel, Message: text, Timestamp: now}
	err = ds.SaveSession(sess)
	if err != nil {
		t.Error(err)
	}

	msg := &Message{SID: sess.ID, Source: tel, Message: text, Timestamp: now}
	err = ds.SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}

	if msg.ID <= 0 {
		t.Error("Failed to set message ID after insert")
	}

	messages, err := ds.FetchAllMessages(sess.ID)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 1 {
		t.Error("Failed to fetch all messages")
	}

	for _, m := range messages {
		if m.Source != tel {
			t.Errorf("Incorrect source: got '%s' should be '%s'", m.Source, tel)
		}
		if m.SID != sess.ID {
			t.Errorf("Incorrect session_id: got '%s' should be '%s'", m.SID, sess.ID)
		}
		if m.Message != text {
			t.Errorf("Incorrect message text: got '%s' should be '%s'", m.Message, text)
		}
	}

	total, err := ds.TotalMessages()
	if err != nil {
		t.Error(err)
	}

	if total != 1 {
		t.Error("Failed to total messages: got '%d' should be '%d'", total, 1)
	}
}

func TestMessageDelete(t *testing.T) {
	ds, err := newTestDataStore()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	text := "Hello World"
	sid := int64(1)

	now := uint64(time.Now().UnixNano() / 1000000)
	msg := &Message{SID: sid, Source: tel, Message: text, Timestamp: now}
	err = ds.SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}

	if msg.ID <= 0 {
		t.Error("Failed to set message ID after insert")
	}

	messages, err := ds.FetchAllMessages(sid)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 1 {
		t.Error("Failed to fetch all messages")
	}

	err = ds.DeleteMessage(msg.ID)
	if err != nil {
		t.Error(err)
	}

	messages, err = ds.FetchAllMessages(sid)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 0 {
		t.Errorf("Incorrect number of messages: got %d should be 0", len(messages))
	}
}

func TestMessageAttachment(t *testing.T) {
	tel := "+1771111006"
	text := "Hello World"
	sid := int64(1)
	now := uint64(time.Now().UnixNano() / 1000000)
	msg := &Message{SID: sid, Source: tel, Message: text, Timestamp: now}

	data := "dummy attachment"

	a := &textsecure.Attachment{
		R:        strings.NewReader(data),
		MimeType: "text/plain",
	}

	dir, err := ioutil.TempDir("", "attachment-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	err = msg.SaveAttachment(dir, a)
	if err != nil {
		t.Error(err)
	}

	data2, err := ioutil.ReadFile(msg.Attachment)
	if err != nil {
		t.Error(err)
	}

	if data != string(data2) {
		t.Error("Failed to write attachment data to file")
	}
}

func TestSentq(t *testing.T) {
	ds, err := newTestDataStore()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	text := "Hello World"
	sid := int64(1)
	now := uint64(time.Now().UnixNano() / 1000000)
	msg := &Message{SID: sid, Source: tel, Message: text, Timestamp: now}

	err = ds.SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}

	err = ds.QueueSent(msg)
	if err != nil {
		t.Error(err)
	}

	messages, err := ds.FetchSentq()
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 1 {
		t.Errorf("Incorrect number of messages in sentq: got %d should be 1", len(messages))
	}

	err = ds.DequeueSent(msg.ID)
	if err != nil {
		t.Error(err)
	}

	messages, err = ds.FetchSentq()
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 0 {
		t.Errorf("Incorrect number of messages in sentq: got %d should be 0", len(messages))
	}
}
