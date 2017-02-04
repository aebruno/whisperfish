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

package store

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher"
	"golang.org/x/crypto/scrypt"
)

const (
	MessageSchema = `
		create table if not exists message 
		(id integer primary key, session_id integer, source text, message string, timestamp integer,
        sent integer default 0, received integer default 0, flags integer default 0, attachment text, 
		mime_type string, has_attachment integer default 0, outgoing integer default 0)
	`
	SentqSchema = `
		create table if not exists sentq
		(message_id integer primary key, timestamp timestamp)
	`
	SessionSchema = `
		create table if not exists session 
		(id integer primary key, source text, message string, timestamp integer,
		 sent integer default 0, received integer default 0, unread integer default 0,
         is_group integer default 0, group_members text, group_id text, group_name text,
		 has_attachment integer default 0)
	`
)

var DS *DataStore

type DataStore struct {
	dbx *sqlx.DB
}

// Create new data store at path. If salt and password are provided the store will
// be encrypted
func NewDataStore(dbPath, saltPath, password string) (*DataStore, error) {
	dsn := dbPath

	if password != "" && saltPath != "" {
		log.Info("Connecting to encrypted data store")
		key, err := getKey(saltPath, password)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to get key")
			return nil, err
		}

		dsn = fmt.Sprintf("%s?_pragma_key=x'%X'&_pragma_cipher_page_size=4096", dbPath, key)
	}

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(SessionSchema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(MessageSchema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(SentqSchema)
	if err != nil {
		return nil, err
	}

	return &DataStore{dbx: db}, nil
}

// Get salt for encrypted database stored at path
func getSalt(path string) ([]byte, error) {
	salt := make([]byte, 8)

	if _, err := os.Stat(path); err == nil {
		salt, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
	} else {
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(path, salt, 0600)
		if err != nil {
			return nil, err
		}
	}

	return salt, nil
}

// Get raw key data for use with sqlcipher
func getKey(saltPath, password string) ([]byte, error) {
	salt, err := getSalt(saltPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get salt")
		return nil, err
	}

	return scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
}

// Encrypt database and closes connection
func (ds *DataStore) Encrypt(path, password string) error {
	key, err := getKey(filepath.Join(filepath.Dir(path), "salt"), password)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get key")
		return err
	}

	query := fmt.Sprintf("ATTACH DATABASE '%s' AS encrypted KEY \"x'%X'\"", path, key)
	_, err = ds.dbx.Exec(query)
	if err != nil {
		return err
	}

	_, err = ds.dbx.Exec("PRAGMA encrypted.cipher_page_size = 4096;")
	if err != nil {
		return err
	}

	_, err = ds.dbx.Exec("SELECT sqlcipher_export('encrypted');")
	if err != nil {
		return err
	}

	_, err = ds.dbx.Exec("DETACH DATABASE encrypted;")
	if err != nil {
		return err
	}

	ds.dbx = nil

	return nil
}

// Decrypt database and closes connection
func (ds *DataStore) Decrypt(path string) error {
	query := fmt.Sprintf("ATTACH DATABASE '%s' AS plaintext KEY '';", path)
	_, err := ds.dbx.Exec(query)
	if err != nil {
		return err
	}

	_, err = ds.dbx.Exec("SELECT sqlcipher_export('plaintext');")
	if err != nil {
		return err
	}

	_, err = ds.dbx.Exec("DETACH DATABASE plaintext;")
	if err != nil {
		return err
	}

	ds.dbx = nil

	return nil
}
