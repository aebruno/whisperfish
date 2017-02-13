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
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func newTestDataStore() (*DataStore, error) {
	saltFile, err := ioutil.TempFile(os.TempDir(), "whispersalt")
	if err != nil {
		return nil, err
	}
	defer os.Remove(saltFile.Name())
	return NewDataStore(":memory:", saltFile.Name(), "secret")
}

func TestSession(t *testing.T) {
	ds, err := newTestDataStore()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	now := uint64(time.Now().UnixNano() / 1000000)
	sess := &Session{Source: tel, Message: "Hello", Timestamp: now}
	sess.Unread = true
	err = ds.SaveSession(sess)
	if err != nil {
		t.Error(err)
	}

	if sess.ID <= 0 {
		t.Error("Failed to set session ID after insert")
	}

	id := sess.ID

	sess, err = ds.FetchSessionBySource(tel)
	if err != nil {
		t.Error(err)
	}

	if sess.Source != tel || sess.ID != id {
		t.Error("Failed to fetch session by tel")
	}

	sess, err = ds.FetchSession(id)
	if err != nil {
		t.Error(err)
	}

	if sess.Source != tel || sess.ID != id {
		t.Error("Failed to fetch session by id")
	}

	sessions, err := ds.FetchAllSessions()
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 1 {
		t.Error("Failed to fetch all sessions")
	}

	total, err := ds.TotalUnread()
	if err != nil {
		t.Error(err)
	}

	if total != 1 {
		t.Error("Failed to fetch all unread sessions")
	}
}

func TestSessionSave(t *testing.T) {
	ds, err := newTestDataStore()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	now := uint64(time.Now().UnixNano() / 1000000)
	sess := &Session{Source: tel, Message: "Hello", Timestamp: now}

	for i := 0; i < 10; i++ {
		sess.Message = fmt.Sprintf("Hello: %d", i)
		err = ds.SaveSession(sess)
		if err != nil {
			t.Error(err)
		}
	}

	sessions, err := ds.FetchAllSessions()
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 1 {
		t.Errorf("Incorrect number of sessions: got %d should be 1", len(sessions))
	}
}

func TestSessionDelete(t *testing.T) {
	ds, err := newTestDataStore()
	if err != nil {
		t.Error(err)
	}

	tel := "+1771111006"
	now := uint64(time.Now().UnixNano() / 1000000)
	sess := &Session{Source: tel, Message: "Hello", Timestamp: now}

	err = ds.SaveSession(sess)
	if err != nil {
		t.Error(err)
	}

	sessions, err := ds.FetchAllSessions()
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 1 {
		t.Errorf("Incorrect number of sessions: got %d should be 1", len(sessions))
	}

	err = ds.DeleteSession(sess.ID)
	if err != nil {
		t.Error(err)
	}

	sessions, err = ds.FetchAllSessions()
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != 0 {
		t.Errorf("Incorrect number of sessions: got %d should be 0", len(sessions))
	}
}
