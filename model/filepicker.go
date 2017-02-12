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
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type FilePicker struct {
	core.QAbstractListModel

	paths      []string
	searchPath string

	_ map[int]*core.QByteArray `property:"roles"`
	_ func()                   `slot:"search"`
	_ func()                   `constructor:"init"`
}

func (f *FilePicker) init() {
	f.SetRoles(map[int]*core.QByteArray{
		RolePath: core.NewQByteArray2("path", len("path")),
	})

	// XXX Limit to Pictures folder for now? Consider making this configurable
	f.searchPath = filepath.Join(core.QStandardPaths_WritableLocation(core.QStandardPaths__HomeLocation), "Pictures")
	f.paths = make([]string, 0)

	f.ConnectData(f.data)
	f.ConnectRowCount(f.rowCount)
	f.ConnectColumnCount(f.columnCount)
	f.ConnectSearch(f.search)
	f.ConnectRoleNames(f.roleNames)
}

func (f *FilePicker) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() < 0 || index.Row() > len(f.paths) {
		return core.NewQVariant()
	}

	p := f.paths[index.Row()]

	switch role {
	case RolePath:
		return core.NewQVariant14(p)
	default:
		return core.NewQVariant()
	}
}

func (f *FilePicker) rowCount(parent *core.QModelIndex) int {
	return len(f.paths)
}

func (f *FilePicker) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (f *FilePicker) roleNames() map[int]*core.QByteArray {
	return f.Roles()
}

func (f *FilePicker) search() {
	log.Infof("Searching for files in path: %s", f.searchPath)
	f.BeginResetModel()
	f.paths = make([]string, 0)
	f.EndResetModel()

	count := 0
	filepath.Walk(f.searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".jpg", ".JPG", ".jpeg", ".JPEG", ".png", ".PNG", ".gif", ".GIF":
			log.Debugf("Found image: %s", path)
			f.BeginInsertRows(core.NewQModelIndex(), len(f.paths), len(f.paths))
			f.paths = append(f.paths, path)
			f.EndInsertRows()
			count++
		}

		return nil
	})

	log.Infof("Found %d files", count)
}
