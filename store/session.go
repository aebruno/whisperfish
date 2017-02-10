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
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
)

type Session struct {
	ID            int64  `db:"id"`
	Source        string `db:"source"`
	IsGroup       bool   `db:"is_group"`
	GroupID       string `db:"group_id"`
	GroupName     string `db:"group_name"`
	Members       string `db:"group_members"`
	Message       string `db:"message"`
	Section       string `db:"-"`
	Timestamp     uint64 `db:"timestamp"`
	Unread        bool   `db:"unread"`
	Sent          bool   `db:"sent"`
	Received      bool   `db:"received"`
	HasAttachment bool   `db:"has_attachment"`
}

// Returns identity of session contact
func (s *Session) ContactIdentity() string {
	id, err := textsecure.ContactIdentityKey(s.Source)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch contact identity")
		return ""
	}

	return fmt.Sprintf("% 0X", id)
}

func (ds *DataStore) FetchSessionBySource(source string) (*Session, error) {
	session := Session{}
	err := ds.dbx.Get(&session, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.unread,
		s.sent,
		s.is_group,
		s.group_id,
		s.group_name,
		s.group_members,
		s.received,
		s.has_attachment
	from 
		session as s
	where s.source = ?`, source)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (ds *DataStore) FetchSessionByGroupID(groupID string) (*Session, error) {
	session := Session{}
	err := ds.dbx.Get(&session, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.unread,
		s.sent,
		s.is_group,
		s.group_id,
		s.group_name,
		s.group_members,
		s.received,
		s.has_attachment
	from 
		session as s
	where s.group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (ds *DataStore) FetchSession(id int64) (*Session, error) {
	session := Session{}
	err := ds.dbx.Get(&session, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.is_group,
		s.group_id,
		s.group_name,
		s.group_members,
		s.unread,
		s.sent,
		s.received,
		s.has_attachment
	from 
		session as s
	where s.id = ?`, id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (ds *DataStore) FetchAllSessions() ([]*Session, error) {
	sessions := []*Session{}
	err := ds.dbx.Select(&sessions, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.is_group,
		s.group_id,
		s.group_name,
		s.group_members,
		s.unread,
		s.sent,
		s.received,
		s.has_attachment
	from 
		session as s
	order by s.timestamp desc`)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (ds *DataStore) SaveSession(session *Session) error {
	cols := []string{"source", "message", "timestamp", "is_group", "group_id", "group_members", "group_name", "unread", "sent", "received", "has_attachment"}
	if session.ID > int64(0) {
		cols = append(cols, "id")
	}

	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	query := "insert or replace into session ("
	query += strings.Join(cols, ",")
	query += ") values (:" + strings.Join(cols, ",:") + ")"

	res, err := tx.NamedExec(query, session)
	if err != nil {
		return err
	}

	session.ID, err = res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Failed to fetch last insert id for session")
	}

	return nil
}

func (ds *DataStore) DeleteSession(id int64) error {
	_, err := ds.dbx.Exec(`delete from session where id = ?`, id)
	if err != nil {
		return err
	}
	_, err = ds.dbx.Exec(`delete from sentq where message_id in (select id from message where session_id = ?)`, id)

	_, err = ds.dbx.Exec(`delete from message where session_id = ?`, id)

	return err
}

func (ds *DataStore) MarkSessionRead(id int64) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = ds.dbx.Exec(`update session set unread = 0 where id = ?`, id)
	if err != nil {
		return err
	}
	return err
}

func (ds *DataStore) TotalUnread() (int, error) {
	var total int
	err := ds.dbx.Get(&total, `select count(*) as total from session where unread = 1`)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (ds *DataStore) MarkSessionSent(id int64, msg string, ts uint64) error {
	tx, err := ds.dbx.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = ds.dbx.Exec(`update session set timestamp = ?, message = ?, unread = 0, sent = 1 where id = ?`, ts, msg, id)
	return err
}
