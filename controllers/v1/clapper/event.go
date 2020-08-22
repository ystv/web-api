package clapper

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/services/clapper/event"
)

// MonthList returns all events for a month.
func MonthList(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Year incorrect, format /yyyy/mm")
	}
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Month incorrect, format /yyyy/mm")
	}
	e, err := event.ListMonth(year, month)
	if err != nil {
		log.Printf("MonthList failed: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, e)
}

// EventGet handles getting all signups and roles for a given event
func EventGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad video ID")
	}
	e, err := event.Get(id)
	if err != nil {
		log.Printf("EventGet failed: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, e)
}

// EventNew handles creating a new event
func EventNew(c echo.Context) error {
	e := event.Event{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	p, err := people.GetToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	eventID, err := event.New(&e, p.UserID)
	if err != nil {
		log.Printf("EventNew failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, eventID)
}

// EventUpdate updates an existing event
func EventUpdate(c echo.Context) error {
	e := event.Event{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	p, err := people.GetToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = event.Update(&e, p.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, err)
		}
		log.Printf("PositionUpdate failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
