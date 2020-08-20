package clapper

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
