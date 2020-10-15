package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
)

// SetCrew handles setting the user ID for a crew object,
// essentially just signing a person up to the position.
//
// @Summary Set crew user by user token
// @Description Uses JWT to set who is doing the crew position
// @ID set-crew-user-token
// @Tags clapper, crews
// @Param crewid path int true "Crew ID"
// @Success 200
// @Router /v1/internal/clapper/crews/{crewid} [put]
func (r *Repos) SetCrew(c echo.Context) error {
	p, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("NewEvent: failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	crewID, err := strconv.Atoi(c.Param("crewid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crew ID")
	}
	err = r.crew.UpdateUser(c.Request().Context(), crewID, p.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// ResetCrew handles setting the crew position back to empty for
// when a user changes their mind.
//
// @Summary Set crew user by user token
// @Description Uses JWT to set who is doing the crew position to empty
// @ID delete-crew-user-token
// @Tags clapper, crews
// @Param crewid path int true "Crew ID"
// @Success 200
// @Router /v1/internal/clapper/crews/{crewid} [delete]
func (r *Repos) ResetCrew(c echo.Context) error {
	_, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("NewEvent: failed to get token: %w", err)
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
