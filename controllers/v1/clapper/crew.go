package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SetCrew handles setting the user ID for a crew object,
// essentially just signing a person up to the position.
//
// @Summary Set crew user by user token
// @Description Uses JWT to set who is doing the crew position
// @ID set-crew-user-token
// @Tags clapper-crews
// @Param eventid path int true "Event ID"
// @Param signupid path int true "Signup ID"
// @Param crewid path int true "Position ID"
// @Success 200
// @Router /v1/internal/clapper/event/{eventid}/{signupid}/{crewid} [put]
func (r *Repos) SetCrew(c echo.Context) error {
	p, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("SetCrew: failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	crewID, err := strconv.Atoi(c.Param("crewid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crew ID")
	}

	err = r.crew.UpdateUserAndVerify(c.Request().Context(), crewID, p.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// ResetCrew handles setting the crew position back to empty for
// when a user changes their mind.
//
// @Summary Reset crew user to blank
// @Description Uses JWT to set who is doing the crew position to empty
// @ID delete-crew-user-token
// @Tags clapper-crews
// @Param eventid path int true "Event ID"
// @Param signupid path int true "Signup ID"
// @Param crewid path int true "Position ID"
// @Success 200
// @Router /v1/internal/clapper/event/{signupid}/{crewid}/reset [put]
func (r *Repos) ResetCrew(c echo.Context) error {
	_, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("ResetCrew: failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	crewID, err := strconv.Atoi(c.Param("crewid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crew ID")
	}

	// get crew object
	_, err = r.crew.Get(c.Request().Context(), crewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "Crew not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// TODO verify user has permission

	err = r.crew.DeleteUser(c.Request().Context(), crewID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// NewCrew handles creating a new crew, this being a single person
//
// @Summary Add a position to a signup sheet as crew
// @Description Creates a new crew object, that being a single person.
// @ID new-crew
// @Tags clapper-crews
// @Accept json
// @Param eventid path int true "Event ID"
// @Param signupid path int true "Signup ID"
// @Param crewid path int true "Position ID"
// @Success 200
// @Router /v1/internal/clapper/event/{eventid}/{signupid}/{positionid} [post]
func (r *Repos) NewCrew(c echo.Context) error {
	signupID, err := strconv.Atoi(c.Param("signupid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid signup ID")
	}

	positionID, err := strconv.Atoi(c.Param("positionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid position ID")
	}

	err = r.crew.New(c.Request().Context(), signupID, positionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to insert crew: ", err)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteCrew handles deleting crew position
//
// @Summary Delete crew
// @Description deletes a crew position by ID.
// @ID delete-crew
// @Tags clapper-crews
// @Param signupid path int true "Event ID"
// @Param signupid path int true "Signup ID"
// @Param signupid path int true "Crew ID"
// @Success 200
// @Router /v1/internal/clapper/{eventid}/{signupid}/{crewid} [delete]
func (r *Repos) DeleteCrew(c echo.Context) error {
	signupID, err := strconv.Atoi(c.Param("crewid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crew ID")
	}

	err = r.crew.Delete(c.Request().Context(), signupID)
	if err != nil {
		err = fmt.Errorf("DeleteCrew failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
