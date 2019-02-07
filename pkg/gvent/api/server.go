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

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmckind/gvent-api/pkg/gvent/datastore"
	"github.com/jmckind/gvent-api/version"
	log "github.com/sirupsen/logrus"
)

// Run will start the web server.
func Run() {
	db := datastore.NewDatabaseConnection()
	if db == nil {
		log.Fatal("database connection unavailable, cowardly refusing to proceed")
	}
	db.Initialize()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"gvent-api": gin.H{"version": version.Version}})
	})
	NewEventHandler(db, router)

	router.Run("0.0.0.0:8000")
}
