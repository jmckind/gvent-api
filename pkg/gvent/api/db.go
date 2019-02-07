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

package gvent

import (
	"errors"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// ErrNotFound indicates a database record could not be located.
var ErrNotFound = errors.New("the requested item could not be found")

func createDoc(doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB("gvent").Table(table).Insert(doc).RunWrite(DBSession)
	if err != nil {
		log.Errorf("Unable to insert doc in table '%s'. %s", table, err.Error())
		return err
	}

	if res.Errors > 0 {
		err = errors.New(res.FirstError)
		log.Errorf("Server error attempting to insert doc in table '%s'. %s", table, err.Error())
		return err
	}
	return nil
}

func deleteDoc(id string, doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB("gvent").Table(table).Get(id).Delete().RunWrite(DBSession)
	if err != nil {
		log.Errorf("Unable to delete doc with id '%s' in table '%s'. %s", id, table, err.Error())
		return err
	}
	if res.Errors > 0 {
		err = errors.New(res.FirstError)
		log.Errorf("Server error attempting to delete doc with id '%d' in table '%s'. %s", id, table, err.Error())
		return err
	}
	return nil
}

func initDB(session *r.Session) error {
	err := r.DBCreate("gvent").Exec(session)
	if err != nil {
		log.Errorf("Unable to create database. %s", err.Error())
	}

	err = r.DB("gvent").TableCreate("event").Exec(session)
	if err != nil {
		log.Errorf("Unable to create table 'event'. %s", err.Error())
		return err
	}

	return nil
}

func getDoc(id string, doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB("gvent").Table(table).Get(id).Run(DBSession)
	if err != nil {
		log.Errorf("Unable to get doc with id '%s' from table '%s'. %s", id, table, err.Error())
		return err
	}

	err = res.One(doc)
	if err == r.ErrEmptyResult {
		log.Debugf("Doc not found with id '%s' from table '%s'. %s", id, table)
		return ErrNotFound
	}
	if err != nil {
		log.Errorf("Unable to parse doc. %s", err.Error())
		return err
	}

	return nil
}

func getDocs(docs interface{}) error {
	table := tableNameFromType(docs)
	res, err := r.DB("gvent").Table(table).Run(DBSession)
	if err != nil {
		log.Errorf("Unable to get docs from table '%s'. %s", table, err.Error())
		return err
	}

	err = res.All(docs)
	if err != nil {
		log.Errorf("Unable to parse docs. %s", err.Error())
		return err
	}

	return nil
}

func updateDoc(id string, doc interface{}) error {
	table := tableNameFromType(doc)
	res, err := r.DB("gvent").Table(table).Get(id).Update(doc).RunWrite(DBSession)
	if err != nil {
		log.Errorf("Unable to replace doc in table '%s'. %s", table, err.Error())
		return err
	}
	if res.Errors > 0 {
		err = errors.New(res.FirstError)
		log.Errorf("Server error attempting to replace doc in table '%s'. %s", table, err.Error())
		return err
	}
	return nil
}

func openDBSession() *r.Session {
	session, err := r.Connect(r.ConnectOpts{
		Address:    "192.168.99.100",
		InitialCap: 10,
		MaxOpen:    10,
	})
	if err != nil {
		log.Errorf("Unable to connect to server. %s", err.Error())
		return nil
	}
	return session
}

func tableNameFromType(t interface{}) string {
	name := strings.ToLower(reflect.TypeOf(t).String())
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
