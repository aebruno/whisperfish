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
	log "github.com/sirupsen/logrus"
	"github.com/aebruno/whisperfish/settings"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type FilePicker struct {
	core.QAbstractListModel

	paths       []string
	nameFilters []string
	settings    *settings.Settings

	_ map[int]*core.QByteArray `property:"roles"`
	_ func()                   `slot:"search"`
	_ func()                   `constructor:"init"`
}

func (f *FilePicker) init() {
	f.settings = settings.NewSettings(nil)
	f.SetRoles(map[int]*core.QByteArray{
		RolePath: core.NewQByteArray2("path", len("path")),
	})

	f.paths = make([]string, 0)
	f.nameFilters = []string{
		"*.jpg",
		"*.JPG",
		"*.jpeg",
		"*.JPEG",
		"*.png",
		"*.PNG",
		"*.gif",
		"*.GIF",
	}

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

func (f *FilePicker) searchPath(path string) {
	var dir = core.NewQDir2(path)
	for _, info := range dir.EntryInfoList(f.nameFilters, core.QDir__AllDirs|core.QDir__NoDot|core.QDir__NoSymLinks|core.QDir__Files, core.QDir__DirsFirst|core.QDir__Time) {
		if info.FileName() == ".." {
			continue
		} else if info.IsDir() {
			f.searchPath(info.FilePath())
		} else if info.IsFile() {
			log.Debugf("Found image: %s", path)
			f.BeginInsertRows(core.NewQModelIndex(), len(f.paths), len(f.paths))
			f.paths = append(f.paths, info.AbsoluteFilePath())
			f.EndInsertRows()
		}
	}
}

func (f *FilePicker) search() {
	f.BeginResetModel()
	f.paths = make([]string, 0)
	f.EndResetModel()

	for _, path := range f.settings.GetStringList("search_paths") {
		log.Infof("Searching for files in path: %s", path)
		f.searchPath(path)
	}
}
