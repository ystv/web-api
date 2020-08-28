package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/clapper"
)

// SignupNew handles a creating a signup sheet
func (r *Repos) SignupNew(c echo.Context) error {
	// Validate event ID
	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad event ID")
	}
	// Bind request json to signup
	s := clapper.Signup{}
	err = c.Bind(&s)
	if err != nil {
		err = fmt.Errorf("SignupNew: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Check event exists
	// TODO we might want to move this inside the service
	e, err := r.event.Get(c.Request().Context(), eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "No event found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Insert new signup sheet
	signupID, err := r.signup.New(c.Request().Context(), e.EventID, s)
	if err != nil {
		err = fmt.Errorf("SignupNew: failed to insert new signup: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, signupID)
}
