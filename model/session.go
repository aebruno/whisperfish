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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/aebruno/whisperfish/store"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type SessionModel struct {
	core.QAbstractListModel

	sessions []*store.Session

	_ map[int]*core.QByteArray        `property:"roles"`
	_ int                             `property:"unread"`
	_ func(index int)                 `slot:"remove"`
	_ func() int                      `slot:"count"`
	_ func()                          `slot:"reload"`
	_ func(sid int64, markRead bool)  `slot:"add"`
	_ func(sid int64)                 `slot:"markRead"`
	_ func(sid int64)                 `slot:"markReceived"`
	_ func(sid int64, message string) `slot:"markSent"`
	_ func()                          `constructor:"init"`
}

// Initialize model
func (model *SessionModel) init() {
	model.sessions = make([]*store.Session, 0)

	model.SetRoles(map[int]*core.QByteArray{
		RoleID:            core.NewQByteArray2("id", len("id")),
		RoleSource:        core.NewQByteArray2("source", len("source")),
		RoleIsGroup:       core.NewQByteArray2("isGroup", len("isGroup")),
		RoleGroupName:     core.NewQByteArray2("groupName", len("groupName")),
		RoleGroupMembers:  core.NewQByteArray2("groupMembers", len("groupMembers")),
		RoleMessage:       core.NewQByteArray2("message", len("message")),
		RoleSection:       core.NewQByteArray2("section", len("section")),
		RoleTimestamp:     core.NewQByteArray2("timestamp", len("timestamp")),
		RoleUnread:        core.NewQByteArray2("unread", len("unread")),
		RoleSent:          core.NewQByteArray2("sent", len("sent")),
		RoleReceived:      core.NewQByteArray2("received", len("received")),
		RoleHasAttachment: core.NewQByteArray2("hasAttachment", len("hasAttachment")),
	})

	// Slots
	model.ConnectRoleNames(model.roleNames)
	model.ConnectData(model.data)
	model.ConnectColumnCount(model.columnCount)
	model.ConnectRowCount(model.rowCount)
	model.ConnectRemove(model.remove)
	model.ConnectReload(model.reload)
	model.ConnectAdd(model.add)
	model.ConnectMarkRead(model.markRead)
	model.ConnectMarkSent(model.markSent)
	model.ConnectMarkReceived(model.markReceived)
	model.ConnectCount(func() int {
		return model.rowCount(nil)
	})
}

// Returns the data stored under the given role for the item referred to by the
// index.
func (model *SessionModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() < 0 || index.Row() > len(model.sessions) {
		return core.NewQVariant()
	}

	session := model.sessions[index.Row()]
	switch role {
	case RoleID:
		return core.NewQVariant9(session.ID)
	case RoleSource:
		return core.NewQVariant14(session.Source)
	case RoleIsGroup:
		return core.NewQVariant11(session.IsGroup)
	case RoleGroupID:
		return core.NewQVariant14(session.GroupID)
	case RoleGroupName:
		return core.NewQVariant14(session.GroupName)
	case RoleGroupMembers:
		return core.NewQVariant14(session.Members)
	case RoleMessage:
		return core.NewQVariant14(session.Message)
	case RoleSection:
		return core.NewQVariant14(session.Section)
	case RoleTimestamp:
		return core.NewQVariant10(session.Timestamp)
	case RoleUnread:
		return core.NewQVariant11(session.Unread)
	case RoleSent:
		return core.NewQVariant11(session.Sent)
	case RoleReceived:
		return core.NewQVariant11(session.Received)
	case RoleHasAttachment:
		return core.NewQVariant11(session.HasAttachment)
	default:
		return core.NewQVariant()
	}
}

// Returns the number of items in the model.
func (model *SessionModel) rowCount(parent *core.QModelIndex) int {
	return len(model.sessions)
}

// Return the number of columns. This will always be 1
func (model *SessionModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

// Return the roles for the model
func (model *SessionModel) roleNames() map[int]*core.QByteArray {
	return model.Roles()
}

// Add or replace a Session in the model.
func (model *SessionModel) add(sid int64, markRead bool) {
	sess, err := store.DS.FetchSession(sid)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   sid,
		}).Error("No session found")
		return
	}

	model.SetSection(sess)

	alreadyUnread := false

	for index, s := range model.sessions {
		if s.ID == sess.ID {
			if s.Unread {
				alreadyUnread = true
			}

			model.BeginRemoveRows(core.NewQModelIndex(), index, index)
			model.sessions = append(model.sessions[:index], model.sessions[index+1:]...)
			model.EndRemoveRows()
			break
		}
	}

	if sess.Unread && markRead {
		store.DS.MarkSessionRead(sid)
		sess.Unread = false
		if alreadyUnread {
			cnt := model.Unread() - 1
			if cnt < 0 {
				cnt = 0
			}
			model.SetUnread(cnt)
			model.UnreadChanged(cnt)
		}
	} else if sess.Unread && !alreadyUnread {
		cnt := model.Unread() + 1
		model.SetUnread(cnt)
		model.UnreadChanged(cnt)
	}

	// Add to top of list
	model.BeginInsertRows(core.NewQModelIndex(), 0, 0)
	model.sessions = append([]*store.Session{sess}, model.sessions...)
	model.EndInsertRows()
}

// Reload all sessions in the model. This clears the model and queries the
// database for the list of sessions
func (model *SessionModel) reload() {
	model.BeginResetModel()
	model.sessions = make([]*store.Session, 0)
	model.EndResetModel()

	sessions, err := store.DS.FetchAllSessions()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to load sessions from database")
		return
	}

	unread := 0
	for _, s := range sessions {
		model.SetSection(s)
		model.BeginInsertRows(core.NewQModelIndex(), len(model.sessions), len(model.sessions))
		model.sessions = append(model.sessions, s)
		model.EndInsertRows()
		if s.Unread {
			unread++
		}
	}
	model.SetUnread(unread)
}

// Mark session as sent
func (model *SessionModel) markSent(id int64, message string) {
	for i, s := range model.sessions {
		if s.ID == id {
			s.Sent = true
			s.Message = message
			var index = model.Index(i, 0, core.NewQModelIndex())
			model.DataChanged(index, index, []int{RoleSent, RoleMessage})
			break
		}
	}
}

// Mark session as received
func (model *SessionModel) markReceived(id int64) {
	for i, s := range model.sessions {
		if s.ID == id {
			s.Received = true
			var index = model.Index(i, 0, core.NewQModelIndex())
			model.DataChanged(index, index, []int{RoleReceived})
			break
		}
	}
}

// Mark session as read
func (model *SessionModel) markRead(id int64) {
	for i, s := range model.sessions {
		if s.ID == id && s.Unread {
			s.Unread = false
			var index = model.Index(i, 0, core.NewQModelIndex())
			model.DataChanged(index, index, []int{RoleUnread})
			store.DS.MarkSessionRead(id)
			cnt := model.Unread() - 1
			if cnt < 0 {
				cnt = 0
			}
			model.SetUnread(cnt)
			model.UnreadChanged(cnt)

			break
		}
	}
}

// Removes session at index. This removes the session from the list model and
// deletes it from the database.
func (model *SessionModel) remove(index int) {
	if index < 0 || index > len(model.sessions)-1 {
		log.WithFields(log.Fields{
			"index": index,
		}).Info("Invalid index for session model")
		return
	}

	session := model.sessions[index]

	err := store.DS.DeleteSession(session.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   session.ID,
		}).Error("Failed to delete session")
	}

	model.BeginRemoveRows(core.NewQModelIndex(), index, index)
	model.sessions = append(model.sessions[:index], model.sessions[index+1:]...)
	model.EndRemoveRows()
}

func (model *SessionModel) SetSection(s *store.Session) {
	ts := time.Unix(0, int64(1000000*s.Timestamp)).Local()
	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	diff := today.Sub(ts)
	if diff.Seconds() <= 0.0 {
		s.Section = core.QCoreApplication_Translate("", "whisperfish-session-section-today", "", -1)
	} else if diff.Seconds() >= 0 && diff.Hours() <= 24 {
		s.Section = core.QCoreApplication_Translate("", "whisperfish-session-section-yesterday", "", -1)
	} else if diff.Seconds() >= 0 && diff.Hours() <= (24*7) {
		dow := ts.Weekday()
		if dow == 0 {
			// In QLocale days are 1 = Monday .. 7 = Sunday
			// In Go days are 0 = Sunday .. 6 = Saturday
			dow = 7
		}
		s.Section = core.QLocale_System().DayName(int(dow), core.QLocale__LongFormat)
	} else {
		s.Section = core.QCoreApplication_Translate("", "whisperfish-session-section-older", "", -1)
	}
}
