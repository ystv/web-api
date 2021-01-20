package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ListTeams handles listing teams and their members and info
//
// @Summary Provides the current teams
// @Description Lists the teams, their members, and info
// @ID get-public-teams
// @Tags public-teams
// @Produce json
// @Success 200 {array} public.Team
// @Router /v1/public/teams [get]
func (r *Repos) ListTeams(c echo.Context) error {
	t, err := r.public.ListTeams(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("Public ListTeams failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, t)
}

// GetTeam handles getting a selected team
//
// @Summary Provides the team of that year
// @Description Contains members and a range of descriptions
// @ID get-public-team
// @Tags public-teams
// @Param teamid path int true "teamid"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/{teamid} [get]
func (r *Repos) GetTeam(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad teamid")
	}
	t, err := r.public.GetTeam(c.Request().Context(), teamID)
	if err != nil {
		err = fmt.Errorf("Public GetTeamByYear failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, t)
}

// GetTeamByYear handles getting teams by calendar year
//
// @Summary Provides the team of a selected year
// @Description Get the team and their members of that year
// @ID get-public-team-year
// @Tags public-teams
// @Param teamid path int true "teamid"
// @Param year path int true "year"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/{teamid}/{year} [get]
func (r *Repos) GetTeamByYear(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad teamid")
	}
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}
	t, err := r.public.GetTeamByYear(c.Request().Context(), teamID, year)
	if err != nil {
		err = fmt.Errorf("Public GetTeamByYear failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, t)
}

// ListOfficers handles listing the current officers
//
// @Summary Provides the current officers
// @Description Lists the current officers, including their info
// @ID get-public-teams-officers-all
// @Tags public-teams
// @Produce json
// @Success 200 {array} public.TeamMember
// @Router /v1/public/teams/officers [get]
func (r *Repos) ListOfficers(c echo.Context) error {
	o, err := r.public.ListOfficers(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("Public GetTeamByYear failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, o)
}
