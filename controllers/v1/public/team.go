package public

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/public"
	"github.com/ystv/web-api/utils"
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
func (s *Store) ListTeams(c echo.Context) error {
	t, err := s.public.ListTeams(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public ListTeams failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(t))
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
func (s *Store) GetTeamByEmail(c echo.Context) error {
	t, err := s.public.GetTeamByEmail(c.Request().Context(), c.Param("emailAlias"))
	if err != nil {
		err = fmt.Errorf("public GetTeamByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByID handles getting a selected team
//
// @Summary Provides the team of that year from id
// @Description Contains members and a range of descriptions
// @ID get-public-team-by-id
// @Tags public-teams
// @Param teamid path int true "teamid"
// @Produce json
// @Success 200 {object} public.Team
// @Router /v1/public/teams/teamid/{teamid} [get]
func (s *Store) GetTeamByID(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad teamid")
	}

	t, err := s.public.GetTeamByID(c.Request().Context(), teamID)
	if err != nil {
		err = fmt.Errorf("public GetTeamByID failed: %w", err)
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
func (s *Store) GetTeamByYearByEmail(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}

	t, err := s.public.GetTeamByYearByEmail(c.Request().Context(), c.Param("emailAlias"), year)
	if err != nil {
		err = fmt.Errorf("public GetTeamByYearByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByYearByID handles getting teams by calendar year
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
func (s *Store) GetTeamByYearByID(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad teamid")
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}

	t, err := s.public.GetTeamByYearByID(c.Request().Context(), teamID, year)
	if err != nil {
		err = fmt.Errorf("public GetTeamByYearByID failed: %w", err)
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
func (s *Store) GetTeamByStartEndYearByEmail(c echo.Context) error {
	startYear, err := strconv.Atoi(c.Param("startYear"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad start year")
	}

	endYear, err := strconv.Atoi(c.Param("endYear"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad end year")
	}

	t, err := s.public.GetTeamByStartEndYearByEmail(c.Request().Context(), c.Param("emailAlias"), startYear, endYear)
	if err != nil {
		err = fmt.Errorf("public GetTeamByStartEndYearByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GetTeamByStartEndYearByID handles getting teams by educational year
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
func (s *Store) GetTeamByStartEndYearByID(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("teamid"))
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

	t, err := s.public.GetTeamByStartEndYearByID(c.Request().Context(), teamID, startYear, endYear)
	if err != nil {
		err = fmt.Errorf("public GetTeamByStartEndYearByID failed: %w", err)
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
func (s *Store) ListOfficers(c echo.Context) error {
	teamMembersDB, err := s.public.ListOfficers(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public GetTeamByYearByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	teamMembers := make([]public.TeamMember, 0)
	for _, m := range teamMembersDB {
		var pronouns *string
		var startDate, endDate *time.Time
		if m.Pronouns.Valid {
			pronouns = &m.Pronouns.String
		}
		if m.StartDate.Valid {
			startDate = &m.StartDate.Time
		}
		if m.EndDate.Valid {
			endDate = &m.EndDate.Time
		}

		teamMembers = append(teamMembers, public.TeamMember{
			UserName:           m.UserName,
			Avatar:             m.Avatar,
			EmailAlias:         m.EmailAlias,
			Pronouns:           pronouns,
			OfficerName:        m.OfficerName,
			OfficerDescription: m.OfficerDescription,
			HistoryWikiURL:     m.HistoryWikiURL,
			TeamEmail:          m.TeamEmail,
			IsLeader:           m.IsLeader,
			IsDeputy:           m.IsDeputy,
			StartDate:          startDate,
			EndDate:            endDate,
		})
	}

	return c.JSON(http.StatusOK, utils.NonNil(teamMembers))
}
