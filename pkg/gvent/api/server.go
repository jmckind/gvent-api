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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmckind/gvent-api/version"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// DBSession holds the global reference to the connection pool.
var DBSession *r.Session

// init will open a connection to the database and initialze.
func init() {
	DBSession = openDBSession()
	initDB(DBSession)
}

// Run will start the web server.
func Run() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"gvent": gin.H{"version": version.Version}})
	})
	addEventRoutes(router)

	router.Run("0.0.0.0:8000")
}
