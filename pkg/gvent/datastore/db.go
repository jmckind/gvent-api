// Copyright 2018 gvent Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datastore

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

const (
	// DefaultHost is the default database hostname.
	DefaultHost = "localhost"

	// DefaultPort is the default database port.
	DefaultPort = "28015"

	// DefaultName is the default database name.
	DefaultName = "gvent"

	// DefaultPoolInitial is the number of connections when the session is created.
	DefaultPoolInitial = 3

	// DefaultPoolMax is the maximum number of connections held in the pool.
	DefaultPoolMax = 10
)

// Database holds the context for the backend datastore.
type Database struct {
	// Name the application database name.
	Name string

	// Session the active database session.
	Session *r.Session
}

// ErrNotFound indicates a database record could not be located.
var ErrNotFound = errors.New("the requested item could not be found")

// CreateDoc will create a new document.
func (db *Database) CreateDoc(doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB(db.Name).Table(table).Insert(doc).RunWrite(db.Session)
	if err != nil {
		log.Errorf("unable to insert doc in table '%s': %s", table, err.Error())
		return err
	}

	if res.Errors > 0 {
		err = errors.New(res.FirstError)
		log.Errorf("server error while inserting doc in table '%s': %s", table, err.Error())
		return err
	}
	return nil
}

// DeleteDoc will delete the document with the given ID.
func (db *Database) DeleteDoc(id string, doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB(db.Name).Table(table).Get(id).Delete().RunWrite(db.Session)
	if err != nil {
		log.Errorf("unable to delete doc with id '%s' in table '%s'. %s", id, table, err.Error())
		return err
	}
	if res.Errors > 0 {
		err = errors.New(res.FirstError)
		log.Errorf("server error while deleting doc with id '%d' in table '%s': %s", id, table, err.Error())
		return err
	}
	return nil
}

// GetDoc will retrieve the document with the given ID.
func (db *Database) GetDoc(id string, doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB(db.Name).Table(table).Get(id).Run(db.Session)
	if err != nil {
		log.Errorf("unable to get doc with id '%s' from table '%s'. %s", id, table, err.Error())
		return err
	}

	err = res.One(doc)
	if err == r.ErrEmptyResult {
		log.Debugf("doc not found with id '%s' in table '%s': %s", id, table)
		return ErrNotFound
	}
	if err != nil {
		log.Errorf("unable to parse doc: %s", err.Error())
		return err
	}

	return nil
}

// GetDocs will retrieve all documents.
func (db *Database) GetDocs(docs interface{}) error {
	table := tableNameFromType(docs)
	res, err := r.DB(db.Name).Table(table).Run(db.Session)
	if err != nil {
		log.Errorf("unable to get docs from table '%s': %s", table, err.Error())
		return err
	}

	err = res.All(docs)
	if err != nil {
		log.Errorf("unable to parse docs: %s", err.Error())
		return err
	}

	return nil
}

// Initialize will create the application database and tables.
func (db *Database) Initialize() error {
	err := r.DBCreate("gvent").Exec(db.Session)
	if err != nil {
		log.Errorf("unable to create database: %s", err.Error())
	}

	err = r.DB(db.Name).TableCreate("event").Exec(db.Session)
	if err != nil {
		log.Errorf("unable to create table 'event': %s", err.Error())
		return err
	}

	return nil
}

// NewDatabaseConnection will return a new database connection.
func NewDatabaseConnection() *Database {
	name := os.Getenv("GVENT_DATABASE_NAME")
	if len(name) <= 0 {
		name = DefaultName
	}
	session, err := r.Connect(buildConnectOpts())
	if err != nil {
		log.Errorf("unable to connect to server: %s", err.Error())
		return nil
	}
	return &Database{Name: name, Session: session}
}

// UpdateDoc will update the given document and ID.
func (db *Database) UpdateDoc(id string, doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB(db.Name).Table(table).Get(id).Update(doc).RunWrite(db.Session)
	if err != nil {
		log.Errorf("unable to update doc in table '%s': %s", table, err.Error())
		return err
	}
	if res.Errors > 0 {
		err = errors.New(res.FirstError)
		log.Errorf("server error while updating doc in table '%s': %s", table, err.Error())
		return err
	}
	return nil
}

// buildConnectOpts will ensure all options have the proper default values.
func buildConnectOpts() r.ConnectOpts {
	host := os.Getenv("GVENT_DATABASE_HOST")
	if len(host) <= 0 {
		host = DefaultHost
	}
	port := os.Getenv("GVENT_DATABASE_PORT")
	if len(port) <= 0 {
		port = DefaultPort
	}

	address := fmt.Sprintf("%s:%s", host, port)
	log.Debugf("database address: %s", address)

	poolInitial, err := strconv.Atoi(os.Getenv("GVENT_DATABASE_POOL_INITIAL"))
	if err != nil {
		poolInitial = DefaultPoolInitial
	}
	log.Debugf("database connection pool initial size: %d", poolInitial)

	poolMax, err := strconv.Atoi(os.Getenv("GVENT_DATABASE_POOL_MAX"))
	if err != nil {
		poolMax = DefaultPoolMax
	}
	log.Debugf("database connection pool max size: %d", poolMax)

	return r.ConnectOpts{
		Address:    address,
		InitialCap: poolInitial,
		MaxOpen:    poolMax,
	}
}

// tableNameFromType will return the name only portion of the Type for the given object.
func tableNameFromType(t interface{}) string {
	name := strings.ToLower(reflect.TypeOf(t).String())
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
