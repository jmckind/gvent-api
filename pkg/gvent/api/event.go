// Copyright [yyyy] [name of copyright owner]
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
	"net/http"

	"github.com/gorilla/mux"
)

func addEventRoutes(r *mux.Router) {
	r.HandleFunc("/", listEvents).Methods("GET").Name("list-events")
	r.HandleFunc("/", newEvent).Methods("POST").Name("new-event")
	r.HandleFunc("/{id}", getEvent).Methods("GET").Name("get-events")
	r.HandleFunc("/{id}", updateEvent).Methods("PUT").Name("update-event")
	r.HandleFunc("/{id}", deleteEvent).Methods("DELETE").Name("delete-event")
}

// listEvents will return a list of Event objects.
func listEvents(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List Events Handler"))
}

// newEvent will create a new Event object.
func newEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("New Event Handler"))
}

// getEvent will retrieve the Event object with the given id.
func getEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Event Handler"))
}

// updateEvent will update the Event object with the given id.
func updateEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update Event Handler"))
}

// deleteEvent will delete the Event object with the given id.
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete Event Handler"))
}
