package main

import (
	"database/sql"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/jmoiron/sqlx"
)

var (
	SESSION_SCHEMA = `
		create table if not exists session 
		(id integer primary key, source text, message string, timestamp timestamp,
		 sent integer default 0, recieved integer default 0, unread integer default 0, is_group integer default 0)
	`
)

type Session struct {
	ID        int64      `db:"id"`
	Source    string     `db:"source"`
	Name      string     `db:"-"`
	IsGroup   bool       `db:"is_group"`
	Message   string     `db:"message"`
	Section   string     `db:"-"`
	Timestamp time.Time  `db:"timestamp"`
	Date      string     `db:"-"`
	Unread    bool       `db:"unread"`
	Sent      bool       `db:"sent"`
	Recieved  bool       `db:"recieved"`
	messages  []*Message `db:"-"`
	Length    int        `db:"-"`
}

type SessionModel struct {
	sessions []*Session
	Length   int
}

func (s *SessionModel) Get(i int) *Session {
	if i == -1 || i >= len(s.sessions) {
		return &Session{}
	}
	return s.sessions[i]
}

func (s *SessionModel) Add(db *sqlx.DB, source string, message string, timestamp time.Time, unread, sent, recieved bool) (*Session, error) {
	sess, err := FetchSessionBySource(db, source)
	if err != nil {
		if err == sql.ErrNoRows {
			sess = &Session{}
		} else {
			return nil, err
		}
	}

	sess.Source = source
	sess.Message = message
	sess.Timestamp = timestamp
	sess.Unread = unread
	sess.Sent = sent
	sess.Recieved = recieved

	err = SaveSession(db, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *SessionModel) Refresh(db *sqlx.DB, c *Contacts) error {
	var err error
	s.sessions, err = FetchAllSessions(db)
	if err != nil {
		return err
	}

	for i := range s.sessions {
		sess := s.sessions[i]
		sess.Name = c.Name(sess.Source)
		sess.UpdateDate()
		s.sessions[i] = sess
	}

	s.Length = len(s.sessions)

	return nil
}

func (s *Session) UpdateDate() {
	ts := s.Timestamp
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	diff := today.Sub(ts)
	if diff.Seconds() <= 0.0 {
		s.Section = "Today"
		s.Date = humanize.RelTime(ts, time.Now(), "", "")
	} else if diff.Seconds() >= 0 && diff.Hours() <= (24*7) {
		s.Section = ts.Weekday().String()
		s.Date = ts.Format("15:04")
	} else {
		s.Section = "Older"
		s.Date = ts.Format("Jan _2")
	}
}

func (s *Session) Refresh(db *sqlx.DB, c *Contacts) {
	var err error
	s.messages, err = FetchAllMessages(db, s.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   s.ID,
		}).Fatal("Failed to fetch messages from database")
	}

	s.Name = c.Name(s.Source)
	s.Length = len(s.messages)
}

func FetchSessionBySource(db *sqlx.DB, source string) (*Session, error) {
	session := Session{}
	err := db.Get(&session, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.unread,
		s.sent,
		s.is_group,
		s.recieved
	from 
		session as s
	where s.source = ?`, source)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func FetchSession(db *sqlx.DB, id int64) (*Session, error) {
	session := Session{}
	err := db.Get(&session, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.is_group,
		s.unread,
		s.sent,
		s.recieved
	from 
		session as s
	where s.id = ?`, id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func FetchAllSessions(db *sqlx.DB) ([]*Session, error) {
	sessions := []*Session{}
	err := db.Select(&sessions, `
	select
		s.id,
		s.source,
		s.message,
		s.timestamp,
		s.is_group,
		s.unread,
		s.sent,
		s.recieved
	from 
		session as s
	order by s.timestamp desc`)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func SaveSession(db *sqlx.DB, session *Session) error {
	cols := []string{"source", "message", "timestamp", "is_group", "unread", "sent", "recieved"}
	if session.ID > int64(0) {
		cols = append(cols, "id")
	}

	query := "insert or replace into session ("
	query += strings.Join(cols, ",")
	query += ") values (:" + strings.Join(cols, ",:") + ")"

	res, err := db.NamedExec(query, session)
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

func DeleteSession(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`delete from session where id = ?`, id)
	if err != nil {
		return err
	}
	_, err = db.Exec(`delete from message where session_id = ?`, id)

	return err
}

func MarkSessionRead(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`update session set unread = 0 where id = ?`, id)
	if err != nil {
		return err
	}
	return err
}
