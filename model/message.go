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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aebruno/textsecure"
	"github.com/aebruno/whisperfish/store"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type MessageModel struct {
	core.QAbstractListModel

	messages []*store.Message

	_ map[int]*core.QByteArray                                            `property:"roles"`
	_ string                                                              `property:"peerIdentity"`
	_ string                                                              `property:"peerName"`
	_ string                                                              `property:"peerTel"`
	_ string                                                              `property:"groupMembers"`
	_ int64                                                               `property:"sessionId"`
	_ bool                                                                `property:"group"`
	_ func(text string)                                                   `slot:"copyToClipboard"`
	_ func() int                                                          `slot:"total"`
	_ func() int                                                          `slot:"unsentCount"`
	_ func(sid int64, peerName string)                                    `slot:"load"`
	_ func(index int)                                                     `slot:"remove"`
	_ func(id int64)                                                      `slot:"add"`
	_ func(id int64)                                                      `slot:"markSent"`
	_ func(id int64)                                                      `slot:"markReceived"`
	_ func(source, message, groupName, attachment string, add bool) int64 `slot:"createMessage"`
	_ func(source string)                                                 `slot:"endSession"`
	_ func()                                                              `slot:"leaveGroup"`
	_ func(index int)                                                     `slot:"openAttachment"`
	_ func(localID, remoteID string)                                      `slot:"addMember"`
	_ func(localID, remoteID string) string                               `slot:"numericFingerprint"`
	_ func(mid int64)                                                     `signal:"sendMessage"`
	_ func()                                                              `constructor:"init"`
}

// Wire up slots
func (model *MessageModel) init() {
	model.messages = make([]*store.Message, 0)

	model.SetRoles(map[int]*core.QByteArray{
		RoleID:            core.NewQByteArray2("id", len("id")),
		RoleSessionID:     core.NewQByteArray2("sid", len("sid")),
		RoleSource:        core.NewQByteArray2("source", len("source")),
		RoleMessage:       core.NewQByteArray2("message", len("message")),
		RoleTimestamp:     core.NewQByteArray2("timestamp", len("timestamp")),
		RoleOutgoing:      core.NewQByteArray2("outgoing", len("outgoing")),
		RoleSent:          core.NewQByteArray2("sent", len("sent")),
		RoleReceived:      core.NewQByteArray2("received", len("received")),
		RoleHasAttachment: core.NewQByteArray2("hasAttachment", len("hasAttachment")),
		RoleAttachment:    core.NewQByteArray2("attachment", len("attachment")),
		RoleMimeType:      core.NewQByteArray2("mimeType", len("mimeType")),
		RoleQueued:        core.NewQByteArray2("queued", len("queued")),
	})

	// Slots
	model.ConnectRoleNames(model.roleNames)
	model.ConnectData(model.data)
	model.ConnectColumnCount(model.columnCount)
	model.ConnectRowCount(model.rowCount)
	model.ConnectLoad(model.load)
	model.ConnectRemove(model.remove)
	model.ConnectAdd(model.add)
	model.ConnectMarkSent(model.markSent)
	model.ConnectMarkReceived(model.markReceived)
	model.ConnectCreateMessage(model.createMessage)
	model.ConnectEndSession(model.endSession)
	model.ConnectLeaveGroup(model.leaveGroup)
	model.ConnectOpenAttachment(model.openAttachment)
	model.ConnectAddMember(model.addMember)
	model.ConnectNumericFingerprint(model.numericFingerprint)

	model.ConnectTotal(func() int {
		total, _ := store.DS.TotalMessages()
		return total
	})
	model.ConnectUnsentCount(func() int {
		messages, _ := store.DS.FetchSentq()
		return len(messages)
	})
}

// Returns the data stored under the given role for the item referred to by the
// index.
func (model *MessageModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() < 0 || index.Row() > len(model.messages) {
		return core.NewQVariant()
	}

	message := model.messages[index.Row()]
	switch role {
	case RoleID:
		return core.NewQVariant9(message.ID)
	case RoleSessionID:
		return core.NewQVariant9(message.SID)
	case RoleSource:
		return core.NewQVariant14(message.Source)
	case RoleMessage:
		return core.NewQVariant14(message.Message)
	case RoleTimestamp:
		return core.NewQVariant10(message.Timestamp)
	case RoleOutgoing:
		return core.NewQVariant11(message.Outgoing)
	case RoleSent:
		return core.NewQVariant11(message.Sent)
	case RoleReceived:
		return core.NewQVariant11(message.Received)
	case RoleHasAttachment:
		return core.NewQVariant11(message.HasAttachment)
	case RoleAttachment:
		return core.NewQVariant14(message.Attachment)
	case RoleMimeType:
		return core.NewQVariant14(message.MimeType)
	case RoleQueued:
		return core.NewQVariant11(message.Queued)
	default:
		return core.NewQVariant()
	}
}

// Returns the number of items in the model.
func (model *MessageModel) rowCount(parent *core.QModelIndex) int {
	return len(model.messages)
}

// Return the number of columns. This will always be 1
func (model *MessageModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

// Return the roles for the model
func (model *MessageModel) roleNames() map[int]*core.QByteArray {
	return model.Roles()
}

// Add MessageObject to list.
func (model *MessageModel) add(id int64) {
	m, err := store.DS.FetchMessage(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Error("No message found")
		return
	}

	model.BeginInsertRows(core.NewQModelIndex(), 0, 0)
	model.messages = append([]*store.Message{m}, model.messages...)
	model.EndInsertRows()
}

// Mark message as sent
func (model *MessageModel) markSent(id int64) {
	model.mark(id, true, false)
}

// Mark message as received
func (model *MessageModel) markReceived(id int64) {
	model.mark(id, false, true)
}

// Mark message as sent/received.
func (model *MessageModel) mark(mid int64, sent, received bool) {
	for i, m := range model.messages {
		if m.ID == mid && m.SID == model.SessionId() {
			var index = model.Index(i, 0, core.NewQModelIndex())
			if sent {
				m.Sent = true
				m.Queued = false
				model.DataChanged(index, index, []int{RoleSent})
				model.DataChanged(index, index, []int{RoleQueued})
			}
			if received {
				m.Received = true
				model.DataChanged(index, index, []int{RoleReceived})
			}
			break
		}
	}
}

// Load all messages for given session id.
func (model *MessageModel) load(sid int64, peerName string) {
	model.BeginResetModel()
	model.messages = make([]*store.Message, 0)
	model.EndResetModel()

	sess, err := store.DS.FetchSession(sid)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   sid,
		}).Error("No session found")
		return
	}

	model.SetSessionId(sid)
	model.SessionIdChanged(sid)

	model.SetGroup(sess.IsGroup)
	model.GroupChanged(sess.IsGroup)

	if sess.IsGroup && sess.GroupName != "" {
		peerName = sess.GroupName
	} else if peerName == "" {
		peerName = sess.Source
	}

	model.SetPeerName(peerName)
	model.PeerNameChanged(peerName)

	model.SetPeerTel(sess.Source)
	model.PeerTelChanged(sess.Source)

	model.SetGroupMembers(sess.Members)
	model.GroupMembersChanged(sess.Members)

	remoteIdentity, err := textsecure.ContactIdentityKey(sess.Source)
	if err == nil {
		identity := fmt.Sprintf("% 0X", remoteIdentity)
		model.SetPeerIdentity(identity)
		model.PeerIdentityChanged(identity)
	}

	messages, err := store.DS.FetchAllMessages(sid)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"sessionID": sid,
		}).Error("Failed to load messages from database")
		return
	}

	for _, m := range messages {
		model.BeginInsertRows(core.NewQModelIndex(), len(model.messages), len(model.messages))
		model.messages = append(model.messages, m)
		model.EndInsertRows()
	}
}

// Removes message at index. This removes the message from the list model and
// deletes it from the database. This should only be called from main QT
// thread.
func (model *MessageModel) remove(index int) {
	if index < 0 || index > len(model.messages)-1 {
		log.WithFields(log.Fields{
			"index": index,
		}).Info("Invalid index for message model")
		return
	}

	message := model.messages[index]

	err := store.DS.DeleteMessage(message.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"sid":   message.ID,
		}).Error("Failed to delete message from database")
	}

	model.BeginRemoveRows(core.NewQModelIndex(), index, index)
	model.messages = append(model.messages[:index], model.messages[index+1:]...)
	model.EndRemoveRows()
}

// Create a new outgoing message, save to database and queue for delivery. If add
// is true then the new message will be appended to the model. When called from
// the NewMessage page, add should be set to false because there is no active
// session. Returns the session ID the message was created under.
func (model *MessageModel) createMessage(source, message, groupName, attachment string, add bool) int64 {
	var group *textsecure.Group

	// If group name or source is a comma separated list then create a group.
	m := strings.Split(source, ",")
	if len(groupName) > 0 || len(m) > 1 {
		var err error
		group, err = textsecure.NewGroup(groupName, m)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"group_name": groupName,
			}).Error("Failed to create new group")
			return 0
		}

		source = group.Hexid
	}

	msg, err := model.queueMessage(source, message, attachment, group)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"source": source,
		}).Error("Failed to add message to queue")
		return 0
	}

	if add {
		model.BeginInsertRows(core.NewQModelIndex(), 0, 0)
		model.messages = append([]*store.Message{msg}, model.messages...)
		model.EndInsertRows()
	}

	model.SendMessage(msg.ID)

	return msg.SID
}

// Perpare outgoing message for delivery to Signal and save message to queue.
// The message will be fetched from the queue by the SendWorker in a separate
// go routine and sent to Signal
func (model *MessageModel) queueMessage(to, msg, attachment string, group *textsecure.Group) (*store.Message, error) {
	message := &store.Message{
		Source:    to,
		Message:   msg,
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Outgoing:  true,
		Queued:    true,
	}

	if len(attachment) > 0 {
		att, err := os.Open(attachment)
		if err != nil {
			return nil, err
		}
		defer att.Close()
		//XXX We have to re-read the attachment to fetch the mime type
		message.MimeType, _ = textsecure.MIMETypeFromReader(att)
		message.Attachment = attachment
		message.HasAttachment = true
	}

	_, err := store.DS.ProcessMessage(message, group, false)
	if err != nil {
		return nil, err
	}

	err = store.DS.QueueSent(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Reset secure session
func (model *MessageModel) endSession(source string) {
	message := &store.Message{
		Source:    source,
		Message:   "[Whisperfish] Reset secure session",
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Outgoing:  true,
		Flags:     textsecure.EndSessionFlag,
	}

	_, err := store.DS.ProcessMessage(message, nil, false)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add EndSession message to database")
		return
	}

	err = store.DS.QueueSent(message)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add EndSession message to queue")
	}

	model.BeginInsertRows(core.NewQModelIndex(), 0, 0)
	model.messages = append([]*store.Message{message}, model.messages...)
	model.EndInsertRows()

	model.SendMessage(message.ID)
}

func (model *MessageModel) xdgOpen(path string) {
	err := exec.Command("xdg-open", path).Run()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  path,
		}).Error("Failed to run xdg-open")
	}
}

// Open attachment
func (model *MessageModel) openAttachment(index int) {
	if index < 0 || index > len(model.messages)-1 {
		log.WithFields(log.Fields{
			"index": index,
		}).Info("Invalid index for message model")
		return
	}

	message := model.messages[index]
	if len(message.Attachment) > 0 {
		go model.xdgOpen(message.Attachment)
	}
}

// Leave group
func (model *MessageModel) leaveGroup() {
	if !model.IsGroup() || len(model.PeerTel()) == 0 {
		return
	}

	err := textsecure.LeaveGroup(model.PeerTel())
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"groupHexID": model.PeerTel(),
			"groupName":  model.PeerName(),
		}).Error("Failed to leave group")
	}
}

// Add group member
func (model *MessageModel) addMember(localID, remoteID string) {
	if !model.IsGroup() || len(model.PeerTel()) == 0 || len(remoteID) == 0 {
		return
	}

	newMembers := make([]string, 0)
	members := strings.Split(model.GroupMembers(), ",")

	for _, m := range members {
		if m == remoteID {
			log.WithFields(log.Fields{
				"groupHexID": model.PeerTel(),
				"groupName":  model.PeerName(),
				"remoteID":   remoteID,
			}).Warn("Already a group member")
			return
		}

		if m != localID {
			newMembers = append(newMembers, m)
		}
	}

	newMembers = append(newMembers, remoteID)

	_, err := textsecure.UpdateGroup(model.PeerTel(), model.PeerName(), newMembers)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"groupHexID": model.PeerTel(),
			"groupName":  model.PeerName(),
			"remoteID":   remoteID,
		}).Error("Failed to add group member")
	}
}

func (model *MessageModel) numericFingerprint(localID, remoteID string) string {
	if model.IsGroup() {
		return ""
	}

	localKey := textsecure.MyIdentityKey()
	remoteKey, err := textsecure.ContactIdentityKey(remoteID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"source": model.PeerTel(),
		}).Error("Failed to fetch contact identity")
		return ""
	}

	fp := store.NumericFingerprint(localID, localKey, remoteID, remoteKey)

	var buffer bytes.Buffer

	for i, n := range fp {
		buffer.WriteRune(n)
		if (i+1)%20 == 0 {
			buffer.WriteString("\n")
		} else if (i+1)%5 == 0 {
			buffer.WriteString("  ")
		}
	}

	return buffer.String()
}
