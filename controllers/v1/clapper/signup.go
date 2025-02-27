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

// NewSignup handles a creating a signup sheet
//
// @Summary New signup sheet
// @Description Creates a new signup sheet; this is the subpart of an event
// @Description containing the list of crew, with a little metadata on top.
// @ID new-signup
// @Tags clapper-signups
// @Accept json
// @Param event body clapper.NewSignup true "Signup object"
// @Success 201 body int "Event ID"
// @Router /v1/internal/clapper/event/{eventid}/signup [post]
func (s *Store) NewSignup(c echo.Context) error {
	// Validate event ID
	eventID, err := strconv.Atoi(c.Param("eventid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad event ID")
	}

	// Bind request json to signup
	var signup clapper.NewSignup

	err = c.Bind(&signup)
	if err != nil {
		err = fmt.Errorf("NewSignup: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Check event exists
	e, err := s.event.Get(c.Request().Context(), eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "No event found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Insert a new signup sheet
	signupID, err := s.signup.New(c.Request().Context(), e.EventID, signup)
	if err != nil {
		err = fmt.Errorf("NewSignup: failed to insert new signup: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, signupID)
}

// UpdateSignup updates an existing signup
//
// @Summary Update signup
// @Description updates a signup sheet, to the body.
// @ID update-signup
// @Tags clapper-signups
// @Param eventid path int true "Event ID"
// @Param signupid path int true "Signup ID"
// @Accept json
// @Param quote body clapper.Signup true "Signup object"
// @Success 200
// @Router /v1/internal/clapper/event/{eventid}/{signupid} [put]
func (s *Store) UpdateSignup(c echo.Context) error {
	var signup clapper.Signup

	err := c.Bind(&signup)
	if err != nil {
		err = fmt.Errorf("UpdateSignup: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	signupID, err := strconv.Atoi(c.Param("signupid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid signup ID")
	}

	signup.SignupID = signupID

	err = s.signup.Update(c.Request().Context(), signup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		err = fmt.Errorf("UpdateSignup failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteSignup handles deleting signup
//
// @Summary Delete signup
// @Description deletes a signup by ID.
// @ID delete-signup
// @Tags clapper-signups
// @Param signupid path int true "Event ID"
// @Param signupid path int true "Signup ID"
// @Success 200
// @Router /v1/internal/clapper/{eventid}/{signupid} [delete]
func (s *Store) DeleteSignup(c echo.Context) error {
	signupID, err := strconv.Atoi(c.Param("signupid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid signup ID")
	}

	err = s.signup.Delete(c.Request().Context(), signupID)
	if err != nil {
		err = fmt.Errorf("DeleteSignup failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
