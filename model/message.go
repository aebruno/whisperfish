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
	log "github.com/Sirupsen/logrus"
	"github.com/aebruno/whisperfish/store"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type MessageObject struct {
	core.QObject

	_ int64  `property:"id"`
	_ int64  `property:"sid"`
	_ string `property:"source"`
	_ string `property:"message"`
	_ uint64 `property:"timestamp"`
	_ bool   `property:"outgoing"`
	_ bool   `property:"sent"`
	_ bool   `property:"received"`
	_ string `property:"attachment"`
	_ string `property:"mimeType"`
	_ bool   `property:"hasAttachment"`
}

//go:generate qtmoc
type MessageModel struct {
	core.QObject

	Model *core.QAbstractListModel
	list  *arraylist.List
	ds    *store.DataStore

	_ string                                                              `property:"peerIdentity"`
	_ string                                                              `property:"peerName"`
	_ string                                                              `property:"peerTel"`
	_ int64                                                               `property:"sessionId"`
	_ bool                                                                `property:"group"`
	_ func(sid int64, source string, group bool)                          `signal:"refresh"`
	_ func(msg *MessageObject)                                            `signal:"update"`
	_ func(mid int64)                                                     `signal:"markSent"`
	_ func(mid int64)                                                     `signal:"markReceived"`
	_ func() int                                                          `slot:"total"`
	_ func() int                                                          `slot:"unsentCount"`
	_ func(msg *MessageObject)                                            `slot:"add"`
	_ func(sid int64, peerName, peerIdentity, peerTel string, group bool) `slot:"load"`
	_ func(mid int64, sent, received bool)                                `slot:"mark"`
	_ func(index int)                                                     `slot:"remove"`
	_ func()                                                              `constructor:"init"`
}

func init() {
	MessageObject_QRegisterMetaType()
}

// Convert message store to QML compatable message QObject
func newMessage(m *store.Message) *MessageObject {
	var msg = NewMessageObject(nil)
	msg.SetId(m.ID)
	msg.SetSid(m.SID)
	msg.SetSource(m.Source)
	msg.SetMessage(m.Message)
	msg.SetTimestamp(m.Timestamp)
	msg.SetOutgoing(m.Outgoing)
	msg.SetSent(m.Sent)
	msg.SetReceived(m.Received)
	msg.SetAttachment(m.Attachment)
	msg.SetMimeType(m.MimeType)
	msg.SetHasAttachment(m.HasAttachment)

	return msg
}

// Dependency inject data store
func (model *MessageModel) SetDataStore(ds *store.DataStore) {
	model.ds = ds
}

// Wire up slots
func (model *MessageModel) init() {
	model.list = arraylist.New()

	model.Model = core.NewQAbstractListModel(nil)
	model.Model.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		return model.data(index, role)
	})
	model.Model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return model.rowCount(parent)
	})
	model.ConnectTotal(func() int {
		total, _ := model.ds.TotalMessages()
		return total
	})
	model.ConnectUnsentCount(func() int {
		messages, _ := model.ds.FetchSentq()
		return len(messages)
	})
	model.ConnectLoad(func(sid int64, peerName, peerIdentity, peerTel string, group bool) {
		model.load(sid, peerName, peerIdentity, peerTel, group)
	})
	model.ConnectAdd(func(msg *MessageObject) {
		model.add(msg)
	})
	model.ConnectRemove(func(index int) {
		model.remove(index)
	})
	model.ConnectMark(func(mid int64, sent, received bool) {
		model.mark(mid, sent, received)
	})
}

// Returns the data stored under the given role for the item referred to by the
// index. This is a required method of the QAbstractListModel. Roles are
// currently unsupported so we just return the entire QObject in the default
// "display" role.
func (model *MessageModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if role != 0 || !index.IsValid() {
		return core.NewQVariant()
	}

	var mp, exists = model.list.Get(index.Row())
	if !exists {
		return core.NewQVariant()
	}

	message := mp.(*MessageObject)
	return message.ToVariant()
}

// Returns the number of items in the list. This is a required method of the
// QAbstractListModel.
func (model *MessageModel) rowCount(parent *core.QModelIndex) int {
	return model.list.Size()
}

// Add MessageObject to list. This can only be called from the main thread.
func (model *MessageModel) add(msg *MessageObject) {

	// msg is a QObject pointer created in a different thread. Before adding
	// to list model we need to create a new pointer from the main QT thread or
	// else qml is unhappy
	var m = NewMessageObject(nil)
	m.SetId(msg.Id())
	m.SetSid(msg.Sid())
	m.SetSource(msg.Source())
	m.SetMessage(msg.Message())
	m.SetTimestamp(msg.Timestamp())
	m.SetOutgoing(msg.IsOutgoing())
	m.SetSent(msg.IsSent())
	m.SetReceived(msg.IsReceived())
	m.SetAttachment(msg.Attachment())
	m.SetMimeType(msg.MimeType())
	m.SetHasAttachment(msg.IsHasAttachment())

	model.Model.BeginInsertRows(core.NewQModelIndex(), 0, 0)
	model.list.Insert(0, m)
	model.Model.EndInsertRows()
}

// Mark message as sent/received. This should only be called from the main
// thread.  This is wired up in QML as follows:
//
//    Connections {
//        target: MessageModel
//        onMarkSent: {
//            MessageModel.mark(mid, true, false)
//        }
//        onMarkReceived: {
//            MessageModel.mark(mid, false, true)
//        }
//    }
func (model *MessageModel) mark(mid int64, sent, received bool) {
	it := model.list.Iterator()
	for it.Next() {
		m := it.Value().(*MessageObject)
		if m.Id() == mid && m.Sid() == model.SessionId() {
			if sent {
				m.SetSent(true)
				m.SentChanged(true)
			}
			if received {
				m.SetReceived(true)
				m.ReceivedChanged(true)
			}
			break
		}
	}
}

// Load all messages for given session id. This should only be called from the
// main thread.  This is wired up in QML as follows:
//
//    Connections {
//        target: MessageModel
//        onRefresh: {
//            MessageModel.load(sid, peerName, peerIdentity, peerTel, group)
//        }
//    }
func (model *MessageModel) load(sid int64, peerName, peerIdentity, peerTel string, group bool) {
	model.Model.BeginResetModel()
	model.list.Clear()
	model.Model.EndResetModel()

	model.SetSessionId(sid)
	model.SetPeerName(peerName)
	model.SetPeerTel(peerTel)
	model.SetPeerIdentity(peerIdentity)
	model.SetGroup(group)
	model.SessionIdChanged(sid)
	model.PeerNameChanged(peerName)
	model.PeerTelChanged(peerTel)
	model.PeerIdentityChanged(peerIdentity)
	model.GroupChanged(group)

	messages, err := model.ds.FetchAllMessages(sid)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"sessionID": sid,
		}).Error("Failed to load messages from database")
		return
	}

	for _, m := range messages {
		model.Model.BeginInsertRows(core.NewQModelIndex(), model.list.Size(), model.list.Size())
		model.list.Add(newMessage(m))
		model.Model.EndInsertRows()
	}
}

// Saves message to database. Should be called from backend thread. No updates
// to the underlying QAbstractListModel are made, instead this method calls the
// Update signal to modify the list model. This is wired up in QML as follows:
//
//    Connections {
//        target: MessageModel
//        onUpdate: {
//            MessageModel.add(msg)
//        }
//    }
//
// Here's the general flow For an incoming message from signal:
//
// 1. New message arrives via websocket from Signal and
//    client.Backend.processMesssage is called
// 2. client.Backend.processMessage calls model.MessageModel.SaveMessage
//    which stores message in database and calls model.MessageModel.Update
//    signal which tells main QT thread to update MessageModel.Model
// 3. In QML we connect MessageModel.Update signal to MessageModel.Add slot
//    which modifies the underlying QAbstractListModel which in turn updates
//    the QML view
//
// The above convoluted process is because we listen for incoming Signal
// messages via websockets in a separate thread which cannot update the
// QAbstractListModel directly. Updates are only allowed from the main thread.
func (model *MessageModel) SaveMessage(message *store.Message) error {
	err := model.ds.SaveMessage(message)
	if err != nil {
		return err
	}

	if model.SessionId() == message.SID {
		model.Update(newMessage(message))
	}

	return nil
}

// Removes message at index. This removes the message from the list model and
// deletes it from the database. This should only be called from main QT
// thread.
func (model *MessageModel) remove(index int) {
	var mp, exists = model.list.Get(index)
	if !exists {
		log.WithFields(log.Fields{
			"index": index,
		}).Info("No message found in model")
		return
	}

	message := mp.(*MessageObject)

	err := model.ds.DeleteMessage(message.Id())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   message.Id(),
		}).Error("Failed to delete message from database")
	}

	model.Model.BeginRemoveRows(core.NewQModelIndex(), index, index)
	model.list.Remove(index)
	model.Model.EndRemoveRows()
}
