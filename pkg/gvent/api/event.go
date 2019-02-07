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
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

// TODO: Remove once DB hooked up!
var events = make(map[string]*Event, 0)

// TODO: Remove once DB hooked up!
func init() {
	sampleEventCount := 3
	for i := 1; i <= sampleEventCount; i++ {
		u4, err := uuid.NewV4()
		if err != nil {
			fmt.Printf("Unable to create UUID: %s", err)
		}

		id := u4.String()
		events[id] = &Event{
			ID:          id,
			Title:       fmt.Sprintf("Test Event %d", i),
			Description: "This is a sample event.",
			StartDate:   time.Now().Format(time.RFC3339),
			EndDate:     time.Now().Format(time.RFC3339),
		}
	}
}

// addEventRoutes will add the routes for Event requests.
func addEventRoutes(e *gin.Engine) {
	e.GET("/events", listEvents)
	e.POST("/events", newEvent)
	e.GET("/events/:id", getEvent)
	e.PUT("/events/:id", updateEvent)
	e.DELETE("/events/:id", deleteEvent)
}

// listEvents will return a list of Event objects.
func listEvents(c *gin.Context) {
	events := make([]Event, 0)
	err := getDocs(&events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
	}
}

// newEvent will create a new Event object.
func newEvent(c *gin.Context) {
	var newEvent Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u4, _ := uuid.NewV4()
	newEvent.ID = u4.String()

	err := createDoc(newEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, gin.H{"event": newEvent})
	}
}

// getEvent will retrieve the Event object with the given id.
func getEvent(c *gin.Context) {
	var event Event
	err := getDoc(c.Param("id"), &event)
	if err == ErrNotFound {
		c.Status(http.StatusNotFound)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"event": event})
	}
}

// updateEvent will update the Event object with the given id.
func updateEvent(c *gin.Context) {
	id := c.Param("id")
	var event Event

	err := getDoc(id, &event)
	if err != nil {
		if err == ErrNotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updEvent Event
	if err := c.ShouldBindJSON(&updEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.Title = updEvent.Title
	event.Description = updEvent.Description
	event.StartDate = updEvent.StartDate
	event.EndDate = updEvent.EndDate

	err = updateDoc(id, event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"event": event})
	}
}

// deleteEvent will delete the Event object with the given id.
func deleteEvent(c *gin.Context) {
	var event Event
	err := getDoc(c.Param("id"), &event)
	if err != nil {
		if err == ErrNotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	err = deleteDoc(c.Param("id"), event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"event": event})
	}
}
