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
	"github.com/gin-gonic/gin"
	"github.com/jmckind/gvent-api/pkg/gvent/datastore"
)

// RequestHandler encapsulates all request handling logic.
type RequestHandler struct {
	DB     *datastore.Database
	Router *gin.Engine
}

// NewRequestHandler will create a new RequestHandler.
func NewRequestHandler(db *datastore.Database, router *gin.Engine) *RequestHandler {
	return &RequestHandler{DB: db, Router: router}
}
