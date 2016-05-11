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

package main

import (
	"github.com/janimo/textsecure"
	"gopkg.in/qml.v1"
)

type DeviceModel struct {
	devices []textsecure.DeviceInfo
	Len      int
}

// Get device by index i
func (d *DeviceModel) Device(i int) textsecure.DeviceInfo {
	if i == -1 {
		return textsecure.DeviceInfo{}
	}
	return d.devices[i]
}

// Refresh list of linked devices
func (d *DeviceModel) Refresh() error {
	linkedDevices, err := textsecure.LinkedDevices()
	if err != nil {
		return err
	}

	d.devices = linkedDevices
	d.Len = len(d.devices)
	qml.Changed(d, &d.Len)

	return nil
}
