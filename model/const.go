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
	"github.com/therecipe/qt/core"
)

// Unique list of Roles for the QAbstractListModels
const (
	RoleID = int(core.Qt__UserRole) + 1<<iota
	RoleSessionID
	RoleSource
	RoleIsGroup
	RoleGroupID
	RoleGroupName
	RoleGroupMembers
	RoleMessage
	RoleSection
	RoleTimestamp
	RoleUnread
	RoleOutgoing
	RoleSent
	RoleReceived
	RoleHasAttachment
	RoleAttachment
	RoleMimeType
	RolePath
	RoleQueued
	RoleName
	RoleCreated
	RoleLastSeen
)
