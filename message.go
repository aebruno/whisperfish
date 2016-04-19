package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

const (
	msgFlagGroupNew     = 1
	msgFlagGroupUpdate  = 2
	msgFlagGroupLeave   = 4
	msgFlagResetSession = 8

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
        sent integer default 0, received integer default 0, flags integer default 0, attachment text, 
		ctype integer default 0)
	`
)

type Message struct {
	ID          int64     `db:"id"`
	SID         int64     `db:"session_id"`
	Source      string    `db:"source"`
	Message     string    `db:"message"`
	Timestamp   time.Time `db:"timestamp"`
	Sent        bool      `db:"sent"`
	Received    bool      `db:"received"`
	Attachment  string    `db:"attachment"`
	Flags       int       `db:"flags"`
	ContentType int       `db:"ctype"`
	Date        string    `db:"date"`
}

type MessageModel struct {
	messages []*Message
	Name     string
	Tel      string
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
	cols := []string{"session_id", "source", "message", "timestamp", "sent", "received", "flags", "attachment", "ctype"}
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
		m.flags,
		m.sent,
		m.received
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

func MarkMessageSent(db *sqlx.DB, id int64, ts time.Time) error {
	_, err := db.Exec(`update message set timestamp = ?, sent = 1 where id = ?`, ts, id)
	return err
}

func MarkMessageReceived(db *sqlx.DB, sid int64, ts time.Time) error {
	_, err := db.Exec(`update message set received = 1 where strftime("%s", timestamp) = strftime("%s", datetime(?, 'unixepoch')) and session_id = ?`, ts.Unix(), sid)
	return err
}
