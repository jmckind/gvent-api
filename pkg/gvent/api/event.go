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
	"github.com/satori/go.uuid"
)

// Event represents an activity at a specified date and time.
type Event struct {
	// ID is the unique identifier for the Event.
	ID string `json:"id" rethinkdb:"id"`

	// Title is a summary for the Event.
	Title string `json:"title" rethinkdb:"title" binding:"required"`

	// Description is the detailed information for the Event.
	Description string `json:"description" rethinkdb:"description" binding:"required"`

	// StartDate is the date and time that the Event starts.
	StartDate string `json:"startDate" rethinkdb:"startDate" binding:"required"`

	// EndDate is the date and time that the Event ends.
	EndDate string `json:"endDate" rethinkdb:"endDate" binding:"required"`
}

// EventHandler encapsulates all event handling logic.
type EventHandler struct {
	RequestHandler
}

// NewEventHandler will create a new EventHandler.
func NewEventHandler(db *datastore.Database, router *gin.Engine) *EventHandler {
	handler := EventHandler{RequestHandler: *NewRequestHandler(db, router)}

	e := router.Group("/events")
	e.GET("/", handler.list)
	e.POST("/", handler.create)
	e.GET("/:id", handler.show)
	e.PUT("/:id", handler.update)
	e.DELETE("/:id", handler.delete)

	return &handler
}

// list will return the list of Events.
func (h *EventHandler) list(c *gin.Context) {
	events := make([]Event, 0)
	err := h.DB.GetDocs(&events)
	if err != nil {
		unexpectedError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
}

// create will create a new Event.
func (h *EventHandler) create(c *gin.Context) {
	var newEvent Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		validationError(c, err)
		return
	}

	newEvent.ID = uuid.NewV4().String()

	err := h.DB.CreateDoc(newEvent)
	if err != nil {
		unexpectedError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"event": newEvent})
}

// show will retrieve the Event with the given ID.
func (h *EventHandler) show(c *gin.Context) {
	var event Event
	err := h.DB.GetDoc(c.Param("id"), &event)
	if err != nil {
		if err == datastore.ErrNotFound {
			c.Status(http.StatusNotFound)
		} else {
			unexpectedError(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

// update will update the Event with the given ID.
func (h *EventHandler) update(c *gin.Context) {
	id := c.Param("id")
	var event Event

	err := h.DB.GetDoc(id, &event)
	if err != nil {
		if err == datastore.ErrNotFound {
			c.Status(http.StatusNotFound)
		} else {
			unexpectedError(c, err)
		}
		return
	}

	var updEvent Event
	if err := c.ShouldBindJSON(&updEvent); err != nil {
		validationError(c, err)
		return
	}

	event.Title = updEvent.Title
	event.Description = updEvent.Description
	event.StartDate = updEvent.StartDate
	event.EndDate = updEvent.EndDate

	err = h.DB.UpdateDoc(id, event)
	if err != nil {
		unexpectedError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

// delete will delete the Event with the given ID.
func (h *EventHandler) delete(c *gin.Context) {
	id := c.Param("id")
	var event Event

	err := h.DB.GetDoc(id, &event)
	if err != nil {
		if err == datastore.ErrNotFound {
			c.Status(http.StatusNotFound)
		} else {
			unexpectedError(c, err)
		}
		return
	}

	err = h.DB.DeleteDoc(id, event)
	if err != nil {
		unexpectedError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
