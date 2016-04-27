package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func newTestDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(SessionSchema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(MessageSchema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestSession(t *testing.T) {
	db, err := newTestDb()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	sess := &Session{Source: tel, Message: "Hello", Timestamp: time.Now()}
	err = SaveSession(db, sess)
	if err != nil {
		t.Error(err)
	}

	if sess.ID <= 0 {
		t.Error("Failed to set session ID after insert")
	}

	id := sess.ID

	sess, err = FetchSessionBySource(db, tel)
	if err != nil {
		t.Error(err)
	}

	if sess.Source != tel || sess.ID != id {
		t.Error("Failed to fetch session by tel")
	}

	sess, err = FetchSession(db, id)
	if err != nil {
		t.Error(err)
	}

	if sess.Source != tel || sess.ID != id {
		t.Error("Failed to fetch session by id")
	}

	sessions, err := FetchAllSessions(db)
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 1 {
		t.Error("Failed to fetch all sessions")
	}
}

func TestSessionSave(t *testing.T) {
	db, err := newTestDb()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	sess := &Session{Source: tel, Message: "Hello", Timestamp: time.Now()}

	for i := 0; i < 10; i++ {
		sess.Message = fmt.Sprintf("Hello: %d", i)
		err = SaveSession(db, sess)
		if err != nil {
			t.Error(err)
		}
	}

	sessions, err := FetchAllSessions(db)
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 1 {
		t.Errorf("Incorrect number of sessions: got %d should be 1", len(sessions))
	}
}

func TestSessionDelete(t *testing.T) {
	db, err := newTestDb()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	sess := &Session{Source: tel, Message: "Hello", Timestamp: time.Now()}

	err = SaveSession(db, sess)
	if err != nil {
		t.Error(err)
	}

	sessions, err := FetchAllSessions(db)
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 1 {
		t.Errorf("Incorrect number of sessions: got %d should be 1", len(sessions))
	}

	err = DeleteSession(db, sess.ID)
	if err != nil {
		t.Error(err)
	}

	sessions, err = FetchAllSessions(db)
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 0 {
		t.Errorf("Incorrect number of sessions: got %d should be 0", len(sessions))
	}
}
