package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

const (
	ContentTypeMessage int = iota
	ContentTypeImage
	ContentTypeVideo
	ContentTypeAudio
)

var (
	MSG_FLAG_RESET_SESSION = 8

	MESSAGE_SCHEMA = `
		create table if not exists message 
		(id integer primary key, session_id integer, source text, message string, timestamp timestamp,
        recieved_at timestamp default null, sent integer default 0, recieved integer default 0,
        flags integer default 0, attachment text, ctype integer default 0)
	`
)

type Message struct {
	ID          int64      `db:"id"`
	SID         int64      `db:"session_id"`
	Source      string     `db:"source"`
	Message     string     `db:"message"`
	Timestamp   time.Time  `db:"timestamp"`
	RecievedAt  *time.Time `db:"recieved_at"`
	Sent        bool       `db:"sent"`
	Recieved    bool       `db:"recieved"`
	Attachment  string     `db:"attachment"`
	Flags       int        `db:"flags"`
	ContentType int        `db:"ctype"`
	Date        string     `db:"date"`
}

type MessageModel struct {
	messages []*Message
	Session  *Session
	Length   int
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
	cols := []string{"session_id", "source", "message", "timestamp", "recieved_at", "sent", "recieved", "flags", "attachment", "ctype"}
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
		m.ctype,
		m.timestamp,
		strftime('%H:%M, %m/%d/%Y', m.timestamp) as date,
		m.recieved_at,
		m.flags,
		m.sent,
		m.recieved
	from 
		message as m
    where m.session_id = ?
	order by m.timestamp asc`, sessionID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func DeleteMessage(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`delete from message where id = ?`, id)
	return err
}
