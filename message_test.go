package main

import (
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	db, err := newTestDb()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	text := "Hello World"
	sess := &Session{Source: tel, Message: text, Timestamp: time.Now()}
	err = SaveSession(db, sess)
	if err != nil {
		t.Error(err)
	}

	msg := &Message{SID: sess.ID, Source: tel, Message: text, Timestamp: time.Now()}
	err = SaveMessage(db, msg)
	if err != nil {
		t.Error(err)
	}

	if msg.ID <= 0 {
		t.Error("Failed to set message ID after insert")
	}

	messages, err := FetchAllMessages(db, sess.ID)
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
}

func TestMessageDelete(t *testing.T) {
	db, err := newTestDb()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	text := "Hello World"
	sid := int64(1)

	msg := &Message{SID: sid, Source: tel, Message: text, Timestamp: time.Now()}
	err = SaveMessage(db, msg)
	if err != nil {
		t.Error(err)
	}

	if msg.ID <= 0 {
		t.Error("Failed to set message ID after insert")
	}

	messages, err := FetchAllMessages(db, sid)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 1 {
		t.Error("Failed to fetch all messages")
	}

	err = DeleteMessage(db, msg.ID)
	if err != nil {
		t.Error(err)
	}

	messages, err = FetchAllMessages(db, sid)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 0 {
		t.Errorf("Incorrect number of messages: got %d should be 0", len(messages))
	}
}
