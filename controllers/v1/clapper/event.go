package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/services/clapper"
)

// MonthList returns all events for a month.
func (r *Repos) MonthList(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Year incorrect, format /yyyy/mm")
	}
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Month incorrect, format /yyyy/mm")
	}
	e, err := r.event.ListMonth(c.Request().Context(), year, month)
	if err != nil {
		err = fmt.Errorf("ListMonth failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, e)
}

// EventGet handles getting all signups and roles for a given event
func (r *Repos) EventGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("eventid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event ID")
	}
	e, err := r.event.Get(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("EventGet failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, e)
}

// EventNew handles creating a new event
func (r *Repos) EventNew(c echo.Context) error {
	e := clapper.Event{}
	err := c.Bind(&e)
	if err != nil {
		err = fmt.Errorf("EventNew: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("EventNew: failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	eventID, err := r.event.New(c.Request().Context(), &e, p.UserID)
	if err != nil {
		err = fmt.Errorf("EventNew: failed to insert event: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, eventID)
}

// EventUpdate updates an existing event
func (r *Repos) EventUpdate(c echo.Context) error {
	e := clapper.Event{}
	err := c.Bind(&e)
	if err != nil {
		err = fmt.Errorf("EventUpdate: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("EventUpdate: failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = r.event.Update(c.Request().Context(), &e, p.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		err = fmt.Errorf("EventUpdate:  failed to update: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// TODO add delete
