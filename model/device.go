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
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/aebruno/textsecure"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type DeviceModel struct {
	core.QAbstractListModel

	devices []textsecure.DeviceInfo

	_ map[int]*core.QByteArray `property:"roles"`
	_ func()                   `constructor:"init"`
	_ func()                   `slot:"reload"`
	_ func(index int)          `slot:"unlink"`
	_ func(tsURL string) bool  `slot:"link"`
}

// Initialize model
func (model *DeviceModel) init() {
	model.devices = make([]textsecure.DeviceInfo, 0)

	model.SetRoles(map[int]*core.QByteArray{
		RoleID:       core.NewQByteArray2("id", len("id")),
		RoleName:     core.NewQByteArray2("name", len("name")),
		RoleCreated:  core.NewQByteArray2("created", len("created")),
		RoleLastSeen: core.NewQByteArray2("lastSeen", len("lastSeen")),
	})

	// Slots
	model.ConnectRoleNames(model.roleNames)
	model.ConnectData(model.data)
	model.ConnectColumnCount(model.columnCount)
	model.ConnectRowCount(model.rowCount)
	model.ConnectReload(model.reload)
	model.ConnectUnlink(model.unlink)
	model.ConnectLink(model.link)
}

// Returns the data stored under the given role for the item referred to by the
// index.
func (model *DeviceModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() < 0 || index.Row() > len(model.devices) {
		return core.NewQVariant()
	}

	device := model.devices[index.Row()]
	switch role {
	case RoleID:
		return core.NewQVariant7(int(device.ID))
	case RoleName:
		return core.NewQVariant14(device.Name)
	case RoleCreated:
		return core.NewQVariant10(device.Created)
	case RoleLastSeen:
		return core.NewQVariant10(device.LastSeen)
	default:
		return core.NewQVariant()
	}
}

// Returns the number of items in the model.
func (model *DeviceModel) rowCount(parent *core.QModelIndex) int {
	return len(model.devices)
}

// Return the number of columns. This will always be 1
func (model *DeviceModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

// Return the roles for the model
func (model *DeviceModel) roleNames() map[int]*core.QByteArray {
	return model.Roles()
}

// Reload all devices in the model. This clears the model and calls
// Signal API for list of devices
func (model *DeviceModel) reload() {
	model.BeginResetModel()
	model.devices = make([]textsecure.DeviceInfo, 0)
	model.EndResetModel()

	linkedDevices, err := textsecure.LinkedDevices()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to fetch linked devices from Signal")
		return
	}

	for _, d := range linkedDevices {
		model.BeginInsertRows(core.NewQModelIndex(), len(model.devices), len(model.devices))
		model.devices = append(model.devices, d)
		model.EndInsertRows()
	}
}

// Unlink a device at index
func (model *DeviceModel) unlink(index int) {
	if index < 0 || index > len(model.devices)-1 {
		log.WithFields(log.Fields{
			"index": index,
		}).Info("Invalid index for device model")
		return
	}

	device := model.devices[index]

	if device.ID == 1 {
		log.Error("Cannot remove the first device")
		return
	}

	err := textsecure.UnlinkDevice(int(device.ID))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to unlink device")
		return
	}

	model.BeginRemoveRows(core.NewQModelIndex(), index, index)
	model.devices = append(model.devices[:index], model.devices[index+1:]...)
	model.EndRemoveRows()
}

// Link a new device with tsURL.
func (model *DeviceModel) link(tsURL string) bool {
	log.WithFields(log.Fields{
		"url": tsURL,
	}).Debug("Linking new device")

	deviceURL, err := url.Parse(tsURL)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to parse URL for new device")
		return false
	}

	uuid := deviceURL.Query().Get("uuid")
	pk := deviceURL.Query().Get("pub_key")
	code, err := textsecure.NewDeviceVerificationCode()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get new device verification code")
		return false
	}

	err = textsecure.AddDevice(uuid, pk, code)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to add device")
		return false
	}

	model.reload()

	return true
}
