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
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type FileObject struct {
	core.QObject

	_ string `property:"path"`
}

//go:generate qtmoc
type FilePicker struct {
	core.QObject

	Model      *core.QAbstractListModel
	list       *arraylist.List
	searchPath string

	_ func() `slot:"search"`
	_ func() `constructor:"init"`
}

func (f *FilePicker) init() {
	f.searchPath = core.QStandardPaths_WritableLocation(core.QStandardPaths__HomeLocation)
	f.list = arraylist.New()

	f.Model = core.NewQAbstractListModel(nil)
	f.Model.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		return f.data(index, role)
	})
	f.Model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return f.rowCount(parent)
	})
	f.ConnectSearch(func() {
		f.search()
	})
}

func (f *FilePicker) data(index *core.QModelIndex, role int) *core.QVariant {
	if role != 0 || !index.IsValid() {
		return core.NewQVariant()
	}

	var ipath, exists = f.list.Get(index.Row())
	if !exists {
		return core.NewQVariant()
	}

	o := ipath.(*FileObject)

	return o.ToVariant()
}

func (f *FilePicker) rowCount(parent *core.QModelIndex) int {
	return f.list.Size()
}

func (f *FilePicker) search() {
	log.Infof("Searching for files in path: %s", f.searchPath)
	f.Model.BeginResetModel()
	f.list.Clear()
	f.Model.EndResetModel()

	count := 0
	filepath.Walk(f.searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".jpg", ".JPG", ".jpeg", ".JPEG", ".png", ".PNG", ".gif", ".GIF":
			log.Debugf("Found image: %s", path)
			f.Model.BeginInsertRows(core.NewQModelIndex(), f.list.Size(), f.list.Size())
			var fo = NewFileObject(nil)
			fo.SetPath(path)
			f.list.Add(fo)
			f.Model.EndInsertRows()
			count++
		}

		return nil
	})

	log.Infof("Found %d files", count)
}
