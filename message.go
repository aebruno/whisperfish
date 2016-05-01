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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/janimo/textsecure"
	"github.com/jmoiron/sqlx"
	"github.com/rogpeppe/fastuuid"
)

const (
	MessageSchema = `
		create table if not exists message 
		(id integer primary key, session_id integer, source text, message string, timestamp timestamp,
        sent integer default 0, received integer default 0, flags integer default 0, attachment text, 
		mime_type string, has_attachment integer default 0, outgoing integer default 0)
	`
	SentqSchema = `
		create table if not exists sentq
		(message_id integer primary key, timestamp timestamp)
	`
)

type Message struct {
	ID            int64     `db:"id"`
	SID           int64     `db:"session_id"`
	Source        string    `db:"source"`
	Message       string    `db:"message"`
	Timestamp     time.Time `db:"timestamp"`
	Outgoing      bool      `db:"outgoing"`
	Sent          bool      `db:"sent"`
	Received      bool      `db:"received"`
	Attachment    string    `db:"attachment"`
	MimeType      string    `db:"mime_type"`
	HasAttachment bool      `db:"has_attachment"`
	Flags         uint32    `db:"flags"`
	Date          string    `db:"date"`
}

type MessageModel struct {
	messages []*Message
	Name     string
	Tel      string
	Length   int
}

func (m *Message) SaveAttachment(dir string, a *textsecure.Attachment) error {
	g, err := fastuuid.NewGenerator()
	if err != nil {
		return err
	}

	uuid := fmt.Sprintf("%x", g.Next())
	adir := filepath.Join(dir, string(uuid[0]))
	os.MkdirAll(adir, 0700)

	fname := filepath.Join(adir, uuid)

	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil
	}
	defer f.Close()

	_, err = io.Copy(f, a.R)
	if err != nil {
		return err
	}

	m.MimeType = a.MimeType
	m.Attachment = fname
	m.HasAttachment = true

	return nil
}

func (m *MessageModel) Get(i int) *Message {
	if i == -1 || i >= len(m.messages) {
		return &Message{}
	}
	return m.messages[i]
}

func (m *MessageModel) RefreshConversation(db *sqlx.DB, sessionID int64) error {
	var err error
	m.messages, err = FetchAllMessages(db, sessionID)
	if err != nil {
		return err
	}

	m.Length = len(m.messages)

	return nil
}

func SaveMessage(db *sqlx.DB, msg *Message) error {
	cols := []string{"session_id", "source", "message", "timestamp", "outgoing", "sent", "received", "flags", "attachment", "mime_type", "has_attachment"}
	if msg.ID > int64(0) {
		cols = append(cols, "id")
	}

	query := "insert or replace into message ("
	query += strings.Join(cols, ",")
	query += ") values (:" + strings.Join(cols, ",:") + ")"

	res, err := db.NamedExec(query, msg)
	if err != nil {
		return err
	}

	msg.ID, err = res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Failed to fetch last insert id for message")
		// XXX Should we bail here?
	}

	return nil
}

func FetchAllMessages(db *sqlx.DB, sessionID int64) ([]*Message, error) {
	messages := []*Message{}
	err := db.Select(&messages, `
	select
		m.id,
		m.session_id,
		m.source,
		m.message,
		m.attachment,
		m.has_attachment,
		m.mime_type,
		m.timestamp,
		strftime('%H:%M, %m/%d/%Y', m.timestamp) as date,
		m.flags,
		m.outgoing,
		m.sent,
		m.received
	from 
		message as m
    where m.session_id = ?
	order by m.timestamp desc`, sessionID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func DeleteMessage(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`delete from message where id = ?`, id)
	return err
}

func MarkMessageSent(db *sqlx.DB, id int64, ts time.Time) error {
	_, err := db.Exec(`update message set timestamp = ?, sent = 1 where id = ?`, ts, id)
	return err
}

func MarkMessageReceived(db *sqlx.DB, source string, ts time.Time) (int64, error) {
	type record struct {
		SessionID int64 `db:"session_id"`
		MessageID int64 `db:"id"`
	}

	rec := record{}

	err := db.Get(&rec, `
		select id,session_id
		from message 
		where strftime("%s", timestamp) = strftime("%s", datetime(?, 'unixepoch'))
		      and sent = 1 and received = 0
	`, ts.Unix())
	if err != nil {
		return rec.SessionID, err
	}

	_, err = db.Exec(`update message set received = 1 where id = ?`, rec.MessageID)
	if err != nil {
		return rec.SessionID, err
	}

	return rec.SessionID, nil
}

func FetchSentq(db *sqlx.DB) ([]*Message, error) {
	messages := []*Message{}
	err := db.Select(&messages, `
	select
		m.id,
		m.session_id,
		m.source,
		m.message,
		m.attachment,
		m.has_attachment,
		m.mime_type,
		m.timestamp,
		strftime('%H:%M, %m/%d/%Y', m.timestamp) as date,
		m.flags,
		m.outgoing,
		m.sent,
		m.received
	from 
		sentq as q
	join message m
	    on q.message_id = m.id
	order by q.timestamp desc`)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func QueueSent(db *sqlx.DB, message *Message) error {
	_, err := db.Exec(`insert into sentq (message_id, timestamp) values (?,datetime('now'))`, message.ID)
	if err != nil {
		return err
	}

	return nil
}

func DequeueSent(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`delete from sentq where message_id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}
