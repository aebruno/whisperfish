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

package model

import (
	"database/sql"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/store"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type SessionObject struct {
	core.QObject

	_ int64  `property:"id"`
	_ string `property:"source"`
	_ string `property:"name"`
	_ bool   `property:"isGroup"`
	_ string `property:"groupId"`
	_ string `property:"groupName"`
	_ string `property:"groupMembers"`
	_ string `property:"message"`
	_ string `property:"section"`
	_ uint64 `property:"timestamp"`
	_ bool   `property:"unread"`
	_ bool   `property:"sent"`
	_ bool   `property:"received"`
	_ bool   `property:"hasAttachment"`
}

//go:generate qtmoc
type SessionModel struct {
	core.QObject

	Model *core.QAbstractListModel
	list  *arraylist.List
	ds    *store.DataStore

	_ int                                        `property:"unread"`
	_ func(sid int64)                            `signal:"markSent"`
	_ func(sid int64)                            `signal:"markReceived"`
	_ func(sid int64)                            `signal:"markRead"`
	_ func()                                     `signal:"refresh"`
	_ func(sess *SessionObject)                  `signal:"update"`
	_ func()                                     `slot:"load"`
	_ func(sess *SessionObject)                  `slot:"add"`
	_ func(sid int64) *SessionObject             `slot:"get"`
	_ func(sid int64, sent, received, read bool) `slot:"mark"`
	_ func(index int)                            `slot:"remove"`
	_ func() int                                 `slot:"count"`
	_ func()                                     `constructor:"init"`
}

func init() {
	SessionObject_QRegisterMetaType()
}

// Convert session store to QML compatable session QObject
func newSession(s *store.Session) *SessionObject {
	s.SetSection()

	var sess = NewSessionObject(nil)
	sess.SetId(s.ID)
	sess.SetSource(s.Source)
	sess.SetName(s.Name)
	sess.SetIsGroup(s.IsGroup)
	sess.SetGroupName(s.GroupName)
	sess.SetGroupMembers(s.Members)
	sess.SetGroupId(s.GroupID)
	sess.SetMessage(s.Message)
	sess.SetSection(s.Section)
	sess.SetTimestamp(s.Timestamp)
	sess.SetUnread(s.Unread)
	sess.SetSent(s.Sent)
	sess.SetReceived(s.Received)
	sess.SetHasAttachment(s.HasAttachment)

	return sess
}

// Dependency inject data store
func (model *SessionModel) SetDataStore(ds *store.DataStore) {
	model.ds = ds
}

// Wire up slots
func (model *SessionModel) init() {
	model.list = arraylist.New()

	model.Model = core.NewQAbstractListModel(nil)
	model.Model.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		return model.data(index, role)
	})
	model.Model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return model.rowCount(parent)
	})
	model.ConnectLoad(func() {
		model.load()
	})
	model.ConnectAdd(func(sess *SessionObject) {
		model.add(sess)
	})
	model.ConnectMark(func(sid int64, sent, received, read bool) {
		model.mark(sid, sent, received, read)
	})
	model.ConnectRemove(func(index int) {
		model.remove(index)
	})
	model.ConnectGet(func(sid int64) *SessionObject {
		return model.get(sid)
	})
	model.ConnectCount(func() int {
		return model.count()
	})
}

// Returns the data stored under the given role for the item referred to by the
// index. This is a required method of the QAbstractListModel. Roles are
// currently unsupported so we just return the entire QObject in the default
// "display" role.
func (model *SessionModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if role != 0 || !index.IsValid() {
		return core.NewQVariant()
	}

	var sp, exists = model.list.Get(index.Row())
	if !exists {
		return core.NewQVariant()
	}

	session := sp.(*SessionObject)
	return session.ToVariant()
}

// Returns the number of items in the list. This is a required method of the
// QAbstractListModel.
func (model *SessionModel) rowCount(parent *core.QModelIndex) int {
	return model.list.Size()
}

// Returns the number sessions. This function is exposed to qml
func (model *SessionModel) count() int {
	return model.list.Size()
}

// Add or replace a SessionObject in the list. This can only be called from the
// main thread.
func (model *SessionModel) add(sess *SessionObject) {
	alreadyUnread := false

	it := model.list.Iterator()
	for it.Next() {
		index, s := it.Index(), it.Value().(*SessionObject)
		if s.Id() == sess.Id() {
			if s.IsUnread() {
				alreadyUnread = true
			}

			// XXX consider moving the row instead of deleting?
			model.Model.BeginRemoveRows(core.NewQModelIndex(), index, index)
			model.list.Remove(index)
			model.Model.EndRemoveRows()
			break
		}
	}

	// sess is a QObject pointer created in a different thread. Before adding
	// to list model we need to create a new pointer from the main QT thread or
	// else qml is unhappy
	var s = NewSessionObject(nil)
	s.SetId(sess.Id())
	s.SetSource(sess.Source())
	s.SetName(sess.Name())
	s.SetIsGroup(sess.IsIsGroup())
	s.SetGroupName(sess.GroupName())
	s.SetGroupMembers(sess.GroupMembers())
	s.SetGroupId(sess.GroupId())
	s.SetMessage(sess.Message())
	s.SetSection(sess.Section())
	s.SetTimestamp(sess.Timestamp())
	s.SetUnread(sess.IsUnread())
	s.SetSent(sess.IsSent())
	s.SetReceived(sess.IsReceived())
	s.SetHasAttachment(sess.IsHasAttachment())

	// Add to top of list
	model.Model.BeginInsertRows(core.NewQModelIndex(), 0, 0)
	model.list.Insert(0, s)
	model.Model.EndInsertRows()

	if sess.IsUnread() && !alreadyUnread {
		cnt := model.Unread() + 1
		model.SetUnread(cnt)
		model.UnreadChanged(cnt)
	}
}

// Load all sessions. This should only be called from the main thread.  This is
// wired up in QML as follows:
//
//    Connections {
//        target: SessionModel
//        onRefresh: {
//            SessionModel.load()
//        }
//    }
func (model *SessionModel) load() {
	model.Model.BeginResetModel()
	model.list.Clear()
	model.Model.EndResetModel()

	sessions, err := model.ds.FetchAllSessions()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to load sessions from database")
		return
	}

	unread := 0
	for _, s := range sessions {
		model.Model.BeginInsertRows(core.NewQModelIndex(), model.list.Size(), model.list.Size())
		model.list.Add(newSession(s))
		model.Model.EndInsertRows()
		if s.Unread {
			unread++
		}
	}
	model.SetUnread(unread)
}

// Mark session as sent, received or read. This should only be called from
// the main thread.  This is wired up in QML as follows:
//
//    Connections {
//        target: SessionModel
//        onMarkSent: {
//            SessionModel.mark(sid, true, false, false)
//        }
//        onMarkReceived: {
//            SessionModel.mark(sid, false, true, false)
//        }
//        onMarkRead: {
//            SessionModel.mark(sid, false, false, true)
//        }
//    }
func (model *SessionModel) mark(id int64, sent, received, read bool) {
	it := model.list.Iterator()
	for it.Next() {
		s := it.Value().(*SessionObject)
		if s.Id() == id {
			if received {
				s.SetReceived(true)
				s.ReceivedChanged(true)
			}
			if sent {
				s.SetSent(true)
				s.SentChanged(true)
			}
			if read && s.IsUnread() {
				s.SetUnread(false)
				s.UnreadChanged(false)
				model.ds.MarkSessionRead(id)
				cnt := model.Unread() - 1
				if cnt < 0 {
					cnt = 0
				}
				model.SetUnread(cnt)
				model.UnreadChanged(cnt)
			}
			break
		}
	}
}

// Process a new message. Create or fetch the session associated with the
// message and save to database. message.SID is set to the corresponding
// session.ID. Should be called from backend thread. No updates to the
// underlying QAbstractListModel are made, instead this method calls the Update
// signal to modify the list model. This is wired up in QML as follows:
//
//    Connections {
//        target: SessionModel
//        onUpdate: {
//            SessionModel.add(sess)
//        }
//    }
//
// Here's the general flow For an incoming message from signal:
//
// 1. New message arrives via websocket from Signal and
//    client.Backend.processMesssage is called
// 2. client.Backend.processMessage calls model.SessionModel.ProcessMessage
//    which stores message in database and calls model.SessionModel.Update
//    signal which tells main QT thread to update SessionModel.Model
// 3. In QML we connect SessionModel.Update signal to SessionModel.Add slot
//    which modifies the underlying QAbstractListModel which in turn updates
//    the QML view
//
// The above convoluted process is because we listen for incoming Signal
// messages via websockets in a separate thread which cannot update the
// QAbstractListModel directly. Updates are only allowed from the main thread.
func (model *SessionModel) ProcessMessage(message *store.Message, group *textsecure.Group, unread bool) error {
	var sess *store.Session
	var err error

	if group != nil {
		sess, err = model.ds.FetchSessionByGroupID(group.Hexid)
	} else {
		sess, err = model.ds.FetchSessionBySource(message.Source)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			sess = &store.Session{}
		} else {
			return err
		}
	}

	if group != nil && group.Flags == textsecure.GroupUpdateFlag {
		message.Message = "Member joined group"
	} else if group != nil && group.Flags == textsecure.GroupLeaveFlag {
		message.Message = "Member left group"
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

	err = model.ds.SaveSession(sess)
	if err != nil {
		return err
	}

	message.SID = sess.ID

	model.Update(newSession(sess))

	return nil
}

// Removes session at index. This removes the session from the list model and
// deletes it from the database. This should only be called from main QT thread.
func (model *SessionModel) remove(index int) {
	var sp, exists = model.list.Get(index)
	if !exists {
		log.WithFields(log.Fields{
			"index": index,
		}).Info("No session found in model")
		return
	}

	session := sp.(*SessionObject)

	err := model.ds.DeleteSession(session.Id())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   session.Id(),
		}).Error("Failed to delete session")
	}

	model.Model.BeginRemoveRows(core.NewQModelIndex(), index, index)
	model.list.Remove(index)
	model.Model.EndRemoveRows()
}

// Get SessionObject with id
func (model *SessionModel) get(sid int64) *SessionObject {
	it := model.list.Iterator()
	for it.Next() {
		s := it.Value().(*SessionObject)
		if s.Id() == sid {
			return s
		}
	}

	return nil
}
