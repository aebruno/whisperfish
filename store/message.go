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
	"database/sql"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/rogpeppe/fastuuid"
)

type Message struct {
	ID            int64  `db:"id"`
	SID           int64  `db:"session_id"`
	Source        string `db:"source"`
	Message       string `db:"message"`
	Timestamp     uint64 `db:"timestamp"`
	Outgoing      bool   `db:"outgoing"`
	Sent          bool   `db:"sent"`
	Received      bool   `db:"received"`
	Attachment    string `db:"attachment"`
	MimeType      string `db:"mime_type"`
	HasAttachment bool   `db:"has_attachment"`
	Flags         uint32 `db:"flags"`
	Queued        bool   `db:"queued"`
}

func (m *Message) SaveAttachment(dir string, a *textsecure.Attachment) error {
	g, err := fastuuid.NewGenerator()
	if err != nil {
		return err
	}

	uuid := fmt.Sprintf("%x", g.Next())
	adir := filepath.Join(dir, string(uuid[0]))
	os.MkdirAll(adir, 0700)

	ext, _ := mime.ExtensionsByType(a.MimeType)
	if ext == nil {
		ext = []string{""}
	}

	fname := filepath.Join(adir, fmt.Sprintf("%s%s", uuid, ext[0]))

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

func (ds *DataStore) SaveMessage(msg *Message) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	cols := []string{"session_id", "source", "message", "timestamp", "outgoing", "sent", "received", "flags", "attachment", "mime_type", "has_attachment"}
	if msg.ID > int64(0) {
		cols = append(cols, "id")
	}

	query := "insert or replace into message ("
	query += strings.Join(cols, ",")
	query += ") values (:" + strings.Join(cols, ",:") + ")"

	res, err := tx.NamedExec(query, msg)
	if err != nil {
		return err
	}

	msg.ID, err = res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Failed to fetch last insert id for message")
		return err
	}

	return nil
}

func (ds *DataStore) FetchAllMessages(sessionID int64) ([]*Message, error) {
	messages := []*Message{}
	err := ds.dbx.Select(&messages, `
	select
		m.id,
		m.session_id,
		m.source,
		m.message,
		m.attachment,
		m.has_attachment,
		m.mime_type,
		m.timestamp,
		m.flags,
		m.outgoing,
		m.sent,
		m.received,
        case when q.message_id > 0 then 1 else 0 end as queued
	from 
		message as m
	left join sentq as q 
        on q.message_id = m.id
    where m.session_id = ?
	order by m.id desc
    `, sessionID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (ds *DataStore) DeleteMessage(id int64) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`delete from message where id = ?`, id)
	return err
}

func (ds *DataStore) DeleteAllMessages(sid int64) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`delete from message where session_id = ?`, sid)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`delete from session where id = ?`, sid)
	return err
}

func (ds *DataStore) MarkMessageSent(id int64, ts uint64) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`update message set timestamp = ?, sent = 1 where id = ?`, ts, id)
	return err
}

func (ds *DataStore) MarkMessageReceived(source string, ts uint64) (int64, int64, error) {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Commit()

	type record struct {
		SessionID int64  `db:"session_id"`
		MessageID int64  `db:"id"`
		Timestamp uint64 `db:"timestamp"`
	}

	rec := record{}

	err = tx.Get(&rec, `
		select id,session_id,timestamp
		from message 
		where timestamp = ?
	`, ts)
	if err != nil {
		return 0, 0, err
	}

	_, err = tx.Exec(`update message set received = 1 where id = ?`, rec.MessageID)
	if err != nil {
		return 0, 0, err
	}

	// only update session if timestamps match
	_, err = tx.Exec(`update session set received = 1 where id = ? and timestamp = ?`, rec.SessionID, rec.Timestamp)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"source":    source,
			"timestamp": rec.Timestamp,
		}).Warn("Failed to mark session received")
		return 0, rec.MessageID, nil
	}

	return rec.SessionID, rec.MessageID, nil
}

func (ds *DataStore) FetchSentq() ([]*Message, error) {
	messages := []*Message{}
	err := ds.dbx.Select(&messages, `
	select
		m.id,
		m.session_id,
		m.source,
		m.message,
		m.attachment,
		m.has_attachment,
		m.mime_type,
		m.timestamp,
		m.flags,
		m.outgoing,
		m.sent,
		m.received,
        1 as queued
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

func (ds *DataStore) FetchQueuedMessage(id int64) (*Message, error) {
	message := Message{}
	err := ds.dbx.Get(&message, `
	select
		m.id,
		m.session_id,
		m.source,
		m.message,
		m.attachment,
		m.has_attachment,
		m.mime_type,
		m.timestamp,
		m.flags,
		m.outgoing,
		m.sent,
		m.received,
        1 as queued
	from 
		sentq as q
	join message m
	    on q.message_id = m.id
    where m.id = ?`, id)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (ds *DataStore) QueueSent(message *Message) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`insert into sentq (message_id, timestamp) values (?,datetime('now'))`, message.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DataStore) DequeueSent(id int64) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec(`delete from sentq where message_id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DataStore) TotalMessages() (int, error) {
	type record struct {
		Total int `db:"total"`
	}

	rec := record{}

	err := ds.dbx.Get(&rec, `select count(*) as total from message`)
	if err != nil {
		return 0, err
	}

	return rec.Total, nil
}

func (ds *DataStore) FetchMessage(id int64) (*Message, error) {
	message := Message{}
	err := ds.dbx.Get(&message, `
	select
		m.id,
		m.session_id,
		m.source,
		m.message,
		m.attachment,
		m.has_attachment,
		m.mime_type,
		m.timestamp,
		m.flags,
		m.outgoing,
		m.sent,
		m.received,
        case when q.message_id > 0 then 1 else 0 end as queued
	from 
		message as m
	left join sentq as q 
        on q.message_id = m.id
    where m.id = ?
    `, id)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// Process message and store in database and update or create a session
func (ds *DataStore) ProcessMessage(message *Message, group *textsecure.Group, unread bool) (*Session, error) {
	var sess *Session
	var err error

	if group != nil {
		sess, err = ds.FetchSessionByGroupID(group.Hexid)
	} else {
		sess, err = ds.FetchSessionBySource(message.Source)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			sess = &Session{}
		} else {
			return nil, err
		}
	}

	sess.Message = message.Message
	sess.Timestamp = message.Timestamp
	sess.Unread = unread
	sess.Sent = message.Sent
	sess.Received = message.Received
	sess.HasAttachment = message.HasAttachment
	if group != nil {
		sess.Source = group.Hexid
		sess.GroupID = group.Hexid
		sess.GroupName = group.Name
		sess.Members = strings.Join(group.Members, ",")
		sess.IsGroup = true
	} else {
		sess.Source = message.Source
	}

	err = ds.SaveSession(sess)
	if err != nil {
		return nil, err
	}

	message.SID = sess.ID
	err = ds.SaveMessage(message)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
