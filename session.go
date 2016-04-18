package main

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/jmoiron/sqlx"
)

var (
	SESSION_SCHEMA = `
		create table if not exists session 
		(id integer primary key, tel text, message string, timestamp timestamp,
		 sent integer default 0, recieved integer default 0, unread integer default 0, is_group integer default 0)
	`
)

type Session struct {
	ID        int64     `db:"id"`
	Tel       string    `db:"tel"`
	Name      string    `db:"-"`
	IsGroup   bool      `db:"is_group"`
	Message   string    `db:"message"`
	Section   string    `db:"-"`
	Timestamp time.Time `db:"timestamp"`
	Date      string    `db:"-"`
	Unread    bool      `db:"unread"`
	Sent      bool      `db:"sent"`
	Recieved  bool      `db:"recieved"`
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

func (s *Session) UpdateDate() {
	ts := s.Timestamp.Local()
	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Local()
	diff := today.Sub(ts)
	if diff.Hours() < 0.0 {
		s.Section = "Today"
		s.Date = humanize.RelTime(ts, time.Now(), "", "")
	} else if diff.Hours() > 0 && diff.Hours() < (24*7) {
		s.Section = ts.Weekday().String()
		s.Date = ts.Format("15:04")
	} else {
		s.Section = "Older"
		s.Date = ts.Format("Jan _2")
	}
}

func (s *SessionModel) AddMessage(db *sqlx.DB, text string, source string, timestamp time.Time) {
	sess, err := FetchSessionByTel(db, source)
	if err != nil {
		if err == sql.ErrNoRows {
			sess = &Session{}
		} else {
			log.WithFields(log.Fields{
				"error": err,
				"tel":   source,
			}).Fatal("Failed to fetch session from database")
			return
		}
	}

	sess.Tel = source
	sess.Message = text
	sess.Timestamp = timestamp

	err = SaveSession(db, sess)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"tel":   source,
		}).Fatal("Failed to update session in database")
		return
	}
}

func (s *SessionModel) Update(db *sqlx.DB, c *Contacts) {
	var err error
	s.sessions, err = FetchAllSessions(db)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to fetch sessions from database")
	}

	for i := range s.sessions {
		sess := s.sessions[i]
		sess.Name = c.Name(sess.Tel)
		sess.UpdateDate()
		s.sessions[i] = sess
	}

	s.Length = len(s.sessions)
}

func FetchSessionByTel(db *sqlx.DB, tel string) (*Session, error) {
	session := Session{}
	err := db.Get(&session, `
	select
		s.id,
		s.tel,
		s.message,
		s.timestamp,
		s.unread,
		s.sent,
		s.is_group,
		s.recieved
	from 
		session as s
	where s.tel = ?`, tel)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func FetchSessionById(db *sqlx.DB, id int64) (*Session, error) {
	session := Session{}
	err := db.Get(&session, `
	select
		s.id,
		s.tel,
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
		s.tel,
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
	query := `insert or replace into session`
	params := []interface{}{
		session.Tel,
		session.Message,
		session.Timestamp,
		session.IsGroup,
		session.Unread,
		session.Sent,
		session.Recieved,
	}
	if session.ID <= int64(0) {
		query += ` (tel,message,timestamp,is_group,unread,sent,recieved) values (?,?,?,?,?,?,?)`
	} else {
		query += ` (id,tel,message,timestamp,is_group,unread,sent,recieved) values (?,?,?,?,?,?,?,?)`
		params = append([]interface{}{session.ID}, params...)
	}

	res, err := db.Exec(query, params...)
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
