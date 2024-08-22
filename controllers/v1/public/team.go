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
		err = fmt.Errorf("public ListTeams failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByEmail handles getting a selected team
//
// @Summary Provides the team of that year from email alias
// @Description Contains members and a range of descriptions
// @ID get-public-team-by-email
// @Tags public-teams
// @Param emailAlias path string true "emailAlias"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/email/{emailAlias} [get]
func (r *Repos) GetTeamByEmail(c echo.Context) error {
	t, err := r.public.GetTeamByEmail(c.Request().Context(), c.Param("emailAlias"))
	if err != nil {
		err = fmt.Errorf("public GetTeamByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamById handles getting a selected team
//
// @Summary Provides the team of that year from id
// @Description Contains members and a range of descriptions
// @ID get-public-team-by-id
// @Tags public-teams
// @Param teamid path int true "teamid"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/teamid/{teamid} [get]
func (r *Repos) GetTeamById(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad teamid")
	}

	t, err := r.public.GetTeamById(c.Request().Context(), teamID)
	if err != nil {
		err = fmt.Errorf("public GetTeamById failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByYearByEmail handles getting teams by calendar year
//
// @Summary Provides the team of a selected year
// @Description Get the team and their members of that year
// @ID get-public-team-year-by-email
// @Tags public-teams
// @Param emailAlias path string true "emailAlias"
// @Param year path int true "year"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/email/{emailAlias}/{year} [get]
func (r *Repos) GetTeamByYearByEmail(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}

	t, err := r.public.GetTeamByYearByEmail(c.Request().Context(), c.Param("emailAlias"), year)
	if err != nil {
		err = fmt.Errorf("public GetTeamByYearByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByYearById handles getting teams by calendar year
//
// @Summary Provides the team of a selected year
// @Description Get the team and their members of that year
// @ID get-public-team-year-by-id
// @Tags public-teams
// @Param teamid path int true "teamid"
// @Param year path int true "year"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/teamid/{teamid}/{year} [get]
func (r *Repos) GetTeamByYearById(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad teamid")
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}

	t, err := r.public.GetTeamByYearById(c.Request().Context(), teamID, year)
	if err != nil {
		err = fmt.Errorf("public GetTeamByYearById failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByStartEndYearByEmail handles getting teams by educational year
//
// @Summary Provides the team of a selected year
// @Description Get the team and their members of that year
// @ID get-public-team-start-end-year-by-email
// @Tags public-teams
// @Param emailAlias path string true "emailAlias"
// @Param startYear path int true "startYear"
// @Param endYear path int true "endYear"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/email/{emailAlias}/{startYear}/{endYear} [get]
func (r *Repos) GetTeamByStartEndYearByEmail(c echo.Context) error {
	startYear, err := strconv.Atoi(c.Param("startYear"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad start year")
	}

	endYear, err := strconv.Atoi(c.Param("endYear"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad end year")
	}

	t, err := r.public.GetTeamByStartEndYearByEmail(c.Request().Context(), c.Param("emailAlias"), startYear, endYear)
	if err != nil {
		err = fmt.Errorf("public GetTeamByStartEndYearByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByStartEndYearById handles getting teams by educational year
//
// @Summary Provides the team of a selected year
// @Description Get the team and their members of that year
// @ID get-public-team-start-end-year-by-id
// @Tags public-teams
// @Param teamid path int true "teamid"
// @Param startYear path int true "startYear"
// @Param endYear path int true "endYear"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/teamid/{teamid}/{startYear}/{endYear} [get]
func (r *Repos) GetTeamByStartEndYearById(c echo.Context) error {
	teamId, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad team id")
	}

	startYear, err := strconv.Atoi(c.Param("startYear"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad start year")
	}

	endYear, err := strconv.Atoi(c.Param("endYear"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad end year")
	}

	t, err := r.public.GetTeamByStartEndYearById(c.Request().Context(), teamId, startYear, endYear)
	if err != nil {
		err = fmt.Errorf("public GetTeamByStartEndYearById failed: %w", err)
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
		err = fmt.Errorf("public GetTeamByYearByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, o)
}
