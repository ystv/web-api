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
// @Summary List events by month
// @Description Lists events by month. The signup section will be null.
// @ID get-events-month
// @Tags events
// @Produce json
// @Param year path int true "year"
// @Param month path int true "month"
// @Success 200 {array} clapper.Event
// @Router /v1/internal/clapper/calendar/{year}/{month} [get]
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
// @Summary Get event by ID
// @Description Get a event including signup-sheets and roles.
// @ID get-event
// @Tags events
// @Produce json
// @Param eventid path int true "Event ID"
// @Success 200 {object} clapper.Event
// @Router /v1/internal/clapper/event/{eventid} [get]
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
// @Summary New event
// @Description creates a new event.
// @Description You do not need to include the sign-up sheets just the meta
// @ID new-event
// @Tags events
// @Accept json
// @Param event body clapper.Event true "Event object"
// @Success 201 body int "Event ID"
// @Router /v1/internal/clapper/event [post]
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
// @Summary New event
// @Description updates an event. Only uses the meta, if you change the
// @Description type it will delete the children.
// @ID update-event
// @Tags events
// @Accept json
// @Param quote body clapper.Event true "Event object"
// @Success 200
// @Router /v1/internal/clapper/event [put]
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
