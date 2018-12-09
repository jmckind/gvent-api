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
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// APIServer holds the context for the running api server.
type APIServer struct {
	router *mux.Router
}

// NewAPIServer returns a new APIServer instance.
func NewAPIServer() *APIServer {
	s := APIServer{}
	s.router = mux.NewRouter()
	return &s
}

// Run will start the web server.
func Run() {
	api := NewAPIServer()
	api.router.HandleFunc("/", IndexHandler)
	api.router.HandleFunc("/events", EventsHandler)

	srv := &http.Server{
		Handler:      api.router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
