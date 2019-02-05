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
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

// addRoutes will populate the router with the routes for the application.
func addRoutes(api *APIServer) {
	api.router.HandleFunc("/", indexHandler)
	addEventRoutes(api.router.PathPrefix("/events").Subrouter())
}

// IndexHandler is responsible for handling requests for the index or home page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Index Handler Called")
}

// Run will start the web server.
func Run() {
	addr := "0.0.0.0:8000"
	api := NewAPIServer()
	addRoutes(api)

	srv := &http.Server{
		Handler:      api.router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Listening at %s", addr)
	log.Fatal(srv.ListenAndServe())
}
