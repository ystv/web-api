package people

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

// OfficershipCount handles getting the officership count
//
// @Summary Provides officership count
// @Description Contains a number of stats
// @ID get-people-officership-count
// @Tags people-officership
// @Produce json
// @Success 200 {object} people.CountOfficerships
// @Router /v1/internal/people/officership/count [get]
func (s *Store) OfficershipCount(c echo.Context) error {
	count, err := s.people.CountOfficerships(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count officerships")
	}

	return c.JSON(http.StatusOK, count)
}

// ListOfficerships handles listing officerships with options
//
// @Summary List officerships with options
// @ID get-people-officerships
// @Tags people-officership
// @Produce json
// @Param officershipStatus query int false "Optional if officership is current (default), retired or any"
// @Success 200 {array} people.Officership
// @Router /v1/internal/people/officerships [get]
func (s *Store) ListOfficerships(c echo.Context) error {
	officershipStatus := c.QueryParam("officershipStatus")

	var dbStatus people.OfficershipsStatus
	switch officershipStatus {
	case "current", "":
		officershipStatus = "current"
		dbStatus = people.Current
	case "retired":
		dbStatus = people.Retired
	case "any":
		dbStatus = people.Any
	default:
		return c.JSON(http.StatusBadRequest,
			errors.New("officershipStatus must be set to either \"any\", \"current\" or \"retired\""))
	}

	officershipsDB, err := s.people.GetOfficerships(c.Request().Context(), dbStatus)
	if err != nil {
		err = fmt.Errorf("ListOfficerships failed to get officerships: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	o := make([]people.Officership, 0)

	for _, officer := range officershipsDB {
		o = append(o, s.people.OfficershipDBToOfficership(officer))
	}

	return c.JSON(http.StatusOK, utils.NonNil(o))
}

// GetOfficership handles getting an officership
//
// @Summary Provides officership
// @Description Contains members and a range of descriptions
// @ID get-people-officership
// @Tags people-officership
// @Param officershipid path int true "officership id"
// @Produce json
// @Success 200 {object} people.Officership
// @Router /v1/internal/people/officership/{officershipid} [get]
func (s *Store) GetOfficership(c echo.Context) error {
	officershipID, err := strconv.Atoi(c.Param("officershipid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid officership id")
	}

	officershipDB, err := s.people.GetOfficership(c.Request().Context(), people.OfficershipGetDTO{OfficershipID: officershipID})
	if err != nil {
		err = fmt.Errorf("GetOfficership failed to get officership: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	o := s.people.OfficershipDBToOfficership(officershipDB)

	return c.JSON(http.StatusOK, o)
}

// AddOfficership handles creating an officership
//
// @Summary Create an officership
// @ID add-people-officership
// @Tags people-officership
// @Produce json
// @Param officership body people.OfficershipAddEditDTO true "Officership object"
// @Success 201 {object} people.Officership
// @Router /v1/internal/people/officership [post]
func (s *Store) AddOfficership(c echo.Context) error {
	var officershipAdd people.OfficershipAddEditDTO
	err := c.Bind(&officershipAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("request body could not be decoded: %w", err))
	}

	if officershipAdd.Name == "" || officershipAdd.EmailAlias == "" || officershipAdd.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name, email alias and description must be filled for add officership")
	}

	if officershipAdd.HistoryWikiURL != "" {
		_, err = url.ParseRequestURI(officershipAdd.HistoryWikiURL)
		if err != nil {
			return fmt.Errorf("failed to parse historyWikiURL: %w", err)
		}
	}

	o1, err := s.people.GetOfficership(c.Request().Context(), people.OfficershipGetDTO{Name: officershipAdd.Name})
	if err == nil && o1.OfficershipID > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "officership with name \""+officershipAdd.Name+"\" already exists")
	}

	officershipDB, err := s.people.AddOfficership(c.Request().Context(), officershipAdd)
	if err != nil {
		return fmt.Errorf("failed to add officerships for add officership: %w", err)
	}

	return c.JSON(http.StatusCreated, s.people.OfficershipDBToOfficership(officershipDB))
}

// EditOfficership handles editing an officership
//
// @Summary Edits an officership
// @ID edit-people-officership
// @Tags people-officership
// @Produce json
// @Param officershipid path int true "officership id"
// @Param officership body people.OfficershipAddEditDTO true "Officership object"
// @Success 200 {object} people.Officership
// @Router /v1/internal/people/officership/{officershipid} [put]
func (s *Store) EditOfficership(c echo.Context) error {
	officershipID, err := strconv.Atoi(c.Param("officershipid"))
	if err != nil {
		return fmt.Errorf("failed to get officershipid for edit officership: %w", err)
	}

	_, err = s.people.GetOfficership(c.Request().Context(),
		people.OfficershipGetDTO{OfficershipID: officershipID})
	if err != nil {
		return fmt.Errorf("failed to get officership for edit officership: %w", err)
	}

	var officershipEdit people.OfficershipAddEditDTO
	err = c.Bind(&officershipEdit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("request body could not be decoded: %w", err))
	}

	if officershipEdit.HistoryWikiURL != "" {
		_, err = url.ParseRequestURI(officershipEdit.HistoryWikiURL)
		if err != nil {
			return fmt.Errorf("failed to parse historyWikiURL: %w", err)
		}
	}

	officershipDB, err := s.people.EditOfficership(c.Request().Context(), officershipID, officershipEdit)
	if err != nil {
		return fmt.Errorf("failed to edit officership for edit officership: %w", err)
	}

	return c.JSON(http.StatusOK, s.people.OfficershipDBToOfficership(officershipDB))
}

// DeleteOfficership handles deleting an officership
//
// @Summary Deletes officership
// @ID delete-people-officership
// @Tags people-officership
// @Param officershipid path int true "officership id"
// @Produce json
// @Success 204
// @Router /v1/internal/people/officership/{officershipid} [delete]
func (s *Store) DeleteOfficership(c echo.Context) error {
	officershipID, err := strconv.Atoi(c.Param("officershipid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to parse officershipid for officership delete: %w", err))
	}

	o, err := s.people.GetOfficership(c.Request().Context(),
		people.OfficershipGetDTO{OfficershipID: officershipID})
	if err != nil {
		return fmt.Errorf("failed to get officership for officership delete: %w", err)
	}

	err = s.people.RemoveOfficershipForOfficershipMembers(c.Request().Context(), officershipID)
	if err != nil {
		return fmt.Errorf("failed to delete officers from officership for officership delete: %w", err)
	}

	if o.TeamID.Valid {
		err = s.people.DeleteOfficershipTeamMember(c.Request().Context(),
			people.OfficershipTeamMemberGetDeleteDTO{OfficerID: officershipID})
		if err != nil {
			return fmt.Errorf("failed to delete team from officership for officership delete: %w", err)
		}
	}

	err = s.people.DeleteOfficership(c.Request().Context(), officershipID)
	if err != nil {
		return fmt.Errorf("failed to delete officership for officership delete: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListOfficershipTeams handles listing officership teams with options
//
// @Summary List officership teams with options
// @ID get-people-officership-teams
// @Tags people-officership-team
// @Produce json
// @Success 200 {array} people.OfficershipTeam
// @Router /v1/internal/people/officership/teams [get]
func (s *Store) ListOfficershipTeams(c echo.Context) error {
	teams, err := s.people.GetOfficershipTeams(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListOfficershipTeams failed to get officership teams: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(teams))
}

// GetOfficershipTeam handles getting an officership team
//
// @Summary Provides officership team
// @Description Contains members and a range of descriptions
// @ID get-people-officership-team
// @Tags people-officership-team
// @Param officershipteamid path int true "officership team id"
// @Produce json
// @Success 200 {object} people.OfficershipTeam
// @Router /v1/internal/people/officership/team/{officershipteamid} [get]
func (s *Store) GetOfficershipTeam(c echo.Context) error {
	officershipTeamId, err := strconv.Atoi(c.Param("officershipteamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid officership team id")
	}

	team, err := s.people.GetOfficershipTeam(c.Request().Context(), people.OfficershipTeamGetDTO{TeamID: officershipTeamId})
	if err != nil {
		err = fmt.Errorf("GetOfficershipTeam failed to get officership team: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, team)
}

// AddOfficershipTeam handles creating an officership team
//
// @Summary Create an officership team
// @ID add-people-officership-team
// @Tags people-officership-team
// @Produce json
// @Param officership body people.OfficershipTeamAddEditDTO true "Officership Team object"
// @Success 201 {object} people.OfficershipTeam
// @Router /v1/internal/people/officership/team [post]
func (s *Store) AddOfficershipTeam(c echo.Context) error {
	var teamAdd people.OfficershipTeamAddEditDTO

	err := c.Bind(&teamAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind add officership team: %w", err))
	}

	t1, err := s.people.GetOfficershipTeam(c.Request().Context(),
		people.OfficershipTeamGetDTO{Name: teamAdd.Name})
	if err == nil && t1.TeamID > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("officership team with name \"%s\" already exists", teamAdd.Name))
	}

	officershipTeam, err := s.people.AddOfficershipTeam(c.Request().Context(), teamAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to add team for addOfficershipTeam: %w", err))
	}

	return c.JSON(http.StatusOK, officershipTeam)
}

// EditOfficershipTeam handles editing an officership team
//
// @Summary Edits an officership team
// @ID edit-people-officership-team
// @Tags people-officership-team
// @Produce json
// @Param officershipteamid path int true "officership team id"
// @Param officership body people.OfficershipTeamAddEditDTO true "Officership Team object"
// @Success 200 {object} people.OfficershipTeam
// @Router /v1/internal/people/officership/team/{officershipteamid} [put]
func (s *Store) EditOfficershipTeam(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("officershipteamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to parse officershipteamid for editOfficershipTeam: %w", err))
	}

	_, err = s.people.GetOfficershipTeam(c.Request().Context(),
		people.OfficershipTeamGetDTO{TeamID: teamID})
	if err != nil {
		return fmt.Errorf("failed to get team for editOfficershipTeam: %w", err)
	}

	var teamEdit people.OfficershipTeamAddEditDTO

	err = c.Bind(&teamEdit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind edit officership team: %w", err))
	}

	officershipTeam, err := s.people.EditOfficershipTeam(c.Request().Context(), teamID, teamEdit)
	if err != nil {
		return fmt.Errorf("failed to edit team for editOfficershipTeam: %w", err)
	}

	return c.JSON(http.StatusOK, officershipTeam)
}

// DeleteOfficershipTeam handles deleting an officership team
//
// @Summary Deletes officership team
// @ID delete-people-officership-team
// @Tags people-officership-team
// @Param officershipteamid path int true "officership team id"
// @Produce json
// @Success 204
// @Router /v1/internal/people/officership/team/{officershipteamid} [delete]
func (s *Store) DeleteOfficershipTeam(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("officershipteamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to parse teamid for officership team delete: %w", err))
	}

	_, err = s.people.GetOfficershipTeam(c.Request().Context(),
		people.OfficershipTeamGetDTO{TeamID: teamID})
	if err != nil {
		return fmt.Errorf("failed to get team for officer team delete: %w", err)
	}

	err = s.people.RemoveTeamForOfficershipTeamMembers(c.Request().Context(), teamID)
	if err != nil {
		return fmt.Errorf("failed to remove officerships from team for officership team delete: %w", err)
	}

	err = s.people.DeleteOfficershipTeam(c.Request().Context(), teamID)
	if err != nil {
		return fmt.Errorf("failed to delete officership team for officership team delete: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListOfficers handles listing officers with options
//
// @Summary List officers with options
// @ID get-people-officers
// @Tags people-officership-officer
// @Produce json
// @Param officershipStatus query int false "Optional if officership is current (default), retired or any"
// @Param officerStatus query int false "Optional if officer is current (default), retired or any"
// @Param officershipId query int false "Optional officership id for getting all officers where id equal"
// @Param userId query int false "Optional user id for getting where user is equal to id"
// @Success 200 {array} people.OfficershipMember
// @Router /v1/internal/people/officership/officers [get]
func (s *Store) ListOfficers(c echo.Context) error {
	officershipStatus := c.QueryParam("officershipStatus")
	officerStatus := c.QueryParam("officerStatus")

	var officershipGet *people.OfficershipGetDTO

	tempOfficershipID, err := strconv.Atoi(c.QueryParam("officershipId"))
	if err == nil {
		officershipGet = &people.OfficershipGetDTO{OfficershipID: tempOfficershipID}
	}

	var userID *int

	tempUserID, err := strconv.Atoi(c.QueryParam("userId"))
	if err == nil {
		userID = &tempUserID
	}

	var dbOfficershipStatus, dbOfficerStatus people.OfficershipsStatus
	switch officershipStatus {
	case "current", "":
		officershipStatus = "current"
		dbOfficershipStatus = people.Current
	case "retired":
		dbOfficershipStatus = people.Retired
	case "any":
		dbOfficershipStatus = people.Any
	default:
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.New("officershipStatus must be set to either \"any\", \"current\" or \"retired\""))
	}

	switch officerStatus {
	case "current", "":
		officerStatus = "current"
		dbOfficerStatus = people.Current
	case "retired":
		dbOfficerStatus = people.Retired
	case "any":
		dbOfficerStatus = people.Any
	default:
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.New("officerStatus must be set to either \"any\", \"current\" or \"retired\""))
	}

	officershipMembersDB, err := s.people.GetOfficershipMembers(c.Request().Context(), officershipGet, userID, dbOfficershipStatus,
		dbOfficerStatus, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get officers: %w", err))
	}

	o := make([]people.OfficershipMember, 0)

	for _, member := range officershipMembersDB {
		o = append(o, s.people.OfficershipMemberDBToOfficershipMember(member))
	}

	return c.JSON(http.StatusOK, utils.NonNil(o))
}

// GetOfficer handles getting an officer
//
// @Summary Provides officer
// @ID get-people-officership-officer
// @Tags people-officership-officer
// @Param officerid path int true "officer id"
// @Produce json
// @Success 200 {object} people.OfficershipMember
// @Router /v1/internal/people/officership/officer/{officerid} [get]
func (s *Store) GetOfficer(c echo.Context) error {
	officerID, err := strconv.Atoi(c.Param("officerid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid officer id")
	}

	officer, err := s.people.GetOfficershipMember(c.Request().Context(), officerID)
	if err != nil {
		return fmt.Errorf("failed to get officer: %w", err)
	}

	return c.JSON(http.StatusOK, s.people.OfficershipMemberDBToOfficershipMember(officer))
}

// AddOfficer handles creating an officer
//
// @Summary Create an officer
// @ID add-people-officership-officer
// @Tags people-officership-officer
// @Produce json
// @Param officership body people.OfficershipMemberAddEditDTO true "Officer object"
// @Success 201 {object} people.OfficershipMember
// @Router /v1/internal/people/officership/officer [post]
func (s *Store) AddOfficer(c echo.Context) error {
	var officerAdd people.OfficershipMemberAddEditDTO

	err := c.Bind(&officerAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to parse officer add data: %w", err))
	}

	if officerAdd.StartDate != nil {
		diffStart := time.Now().Compare(*officerAdd.StartDate)
		if diffStart != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "start date must be before today")
		}

		// Add 22 hours to always be at the end of the day when adding vs the midnight for ending,
		// this takes into consideration daylight savings from the server side
		tempTime := officerAdd.StartDate.Add(time.Hour * 22)
		officerAdd.StartDate = &tempTime
	}

	if officerAdd.EndDate != nil {
		diffEnd := time.Now().Compare(*officerAdd.EndDate)
		if diffEnd != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "end date must be before today")
		}

		if officerAdd.StartDate != nil {
			diffStartEnd := officerAdd.StartDate.Compare(*officerAdd.EndDate)
			if diffStartEnd == 1 {
				return echo.NewHTTPError(http.StatusBadRequest, "start date must be before end date")
			}
		}
	}

	_, err = s.people.GetUser(c.Request().Context(), officerAdd.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get user for officerAdd: %w", err))
	}

	_, err = s.people.GetOfficership(c.Request().Context(), people.OfficershipGetDTO{OfficershipID: officerAdd.OfficerID})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get officership for officerAdd: %w", err))
	}

	officershipMemberDB, err := s.people.AddOfficershipMember(c.Request().Context(), officerAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to add officer for officerAdd: %w", err))
	}

	return c.JSON(http.StatusOK, s.people.OfficershipMemberDBToOfficershipMember(officershipMemberDB))
}

// EditOfficer handles editing an officer
//
// @Summary Edits an officer
// @ID edit-people-officership-officer
// @Tags people-officership-officer
// @Produce json
// @Param officerid path int true "officer id"
// @Param officership body people.OfficershipMemberAddEditDTO true "Officer object"
// @Success 200 {object} people.OfficershipMember
// @Router /v1/internal/people/officership/officer/{officerid} [put]
func (s *Store) EditOfficer(c echo.Context) error {
	officerID, err := strconv.Atoi(c.Param("officerid"))
	if err != nil {
		return fmt.Errorf("failed to get officerid for edit officer: %w", err)
	}

	_, err = s.people.GetOfficershipMember(c.Request().Context(), officerID)
	if err != nil {
		return fmt.Errorf("failed to get officer for edit officer: %w", err)
	}

	userID, err := strconv.Atoi(c.FormValue("userID"))
	if err != nil {
		return fmt.Errorf("failed to get userID form for edit officer: %w", err)
	}

	_, err = s.people.GetUser(c.Request().Context(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user form for edit officer: %w", err)
	}

	officershipID, err := strconv.Atoi(c.FormValue("officershipID"))
	if err != nil {
		return fmt.Errorf("failed to get officershipID form for edit officer: %w", err)
	}

	_, err = s.people.GetOfficership(c.Request().Context(), people.OfficershipGetDTO{OfficershipID: officershipID})
	if err != nil {
		return fmt.Errorf("failed to get officership for edit officer: %w", err)
	}

	tempStartDate := c.FormValue("startDate")
	tempEndDate := c.FormValue("endDate")

	if tempStartDate == "" {
		return errors.New("start date cannot be blank")
	}

	parsedStart, err := time.Parse("02/01/2006", tempStartDate)
	if err != nil {
		return fmt.Errorf("failed to parse start date: %w", err)
	}

	diff := time.Now().Compare(parsedStart)
	if diff != 1 {
		return errors.New("start date must be before today")
	}

	var endDate *time.Time

	if tempEndDate != "" {
		var parsedEnd time.Time

		parsedEnd, err = time.Parse("02/01/2006", tempEndDate)
		if err != nil {
			return fmt.Errorf("failed to parse end date: %w", err)
		}

		endDate = &parsedEnd
	}

	officerEdit := people.OfficershipMemberAddEditDTO{
		UserID:    userID,
		OfficerID: officershipID,
		StartDate: &parsedStart,
		EndDate:   endDate,
	}

	_, err = s.people.EditOfficershipMember(c.Request().Context(), officerID, officerEdit)
	if err != nil {
		return fmt.Errorf("failed to edit officer for edit officer: %w", err)
	}

	officerDB, err := s.people.GetOfficershipMember(c.Request().Context(), officerID)
	if err != nil {
		return fmt.Errorf("failed to get officer for edit officer: %w", err)
	}

	return c.JSON(http.StatusOK, s.people.OfficershipMemberDBToOfficershipMember(officerDB))
}

// DeleteOfficer handles deleting an officer
//
// @Summary Deletes officer
// @ID delete-people-officership-officer
// @Tags people-officership-officer
// @Param officerid path int true "officer id"
// @Produce json
// @Success 204
// @Router /v1/internal/people/officership/officer/{officerid} [delete]
func (s *Store) DeleteOfficer(c echo.Context) error {
	officerID, err := strconv.Atoi(c.Param("officerid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("failed to parse officerid for officer delete: %w", err))
	}

	officer, err := s.people.GetOfficershipMember(c.Request().Context(), officerID)
	if err != nil {
		return fmt.Errorf("failed to get officer for officer delete: %w", err)
	}

	err = s.people.DeleteOfficershipMember(c.Request().Context(), officer.OfficershipMemberID)
	if err != nil {
		return fmt.Errorf("failed to delete officer for officer delete: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}

// OfficershipTeamAddOfficership handles adding an officership to a team
//
// @Summary Create an officership link to a team
// @ID add-people-officership-team-officership
// @Tags people-officership-team
// @Param officershipteamid path int true "officership team id"
// @Produce json
// @Param officership body people.OfficershipTeamMemberAddDTO true "Officership team member object"
// @Success 201 {object} people.OfficershipTeamMember
// @Router /v1/internal/people/officership/team/{officershipteamid}/officership [post]
func (s *Store) OfficershipTeamAddOfficership(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("officershipteamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get team id for officership team add officership: %w", err))
	}

	_, err = s.people.GetOfficershipTeam(c.Request().Context(), people.OfficershipTeamGetDTO{TeamID: teamID})
	if err != nil {
		return fmt.Errorf("failed to get team for officership team add officership: %w", err)
	}

	var teamMemberAdd people.OfficershipTeamMemberAddDTO

	err = c.Bind(&teamMemberAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind data: %w", err))
	}

	_, err = s.people.GetOfficership(c.Request().Context(), people.OfficershipGetDTO{OfficershipID: teamMemberAdd.OfficerID})
	if err != nil {
		return fmt.Errorf("failed to get officership for officership team add officership: %w", err)
	}

	_, err = s.people.GetOfficershipTeamMember(c.Request().Context(), people.OfficershipTeamMemberGetDeleteDTO{
		TeamID:    teamID,
		OfficerID: teamMemberAdd.OfficerID,
	})
	if err == nil {
		return errors.New("failed to add officership team member for officership team add officership: row already exists")
	}

	teamMember, err := s.people.AddOfficershipTeamMember(c.Request().Context(), teamMemberAdd)
	if err != nil {
		return fmt.Errorf("failed to add officership team member for officership team add officership: %w", err)
	}

	return c.JSON(http.StatusOK, teamMember)
}

// OfficershipTeamRemoveOfficership handles deleting an officership from a team
//
// @Summary Removes officership from team
// @ID remove-people-officership-from-team
// @Tags people-officership-team
// @Param officershipteamid path int true "officership team id"
// @Param officershipid path int true "officership id"
// @Produce json
// @Success 204
// @Router /v1/internal/people/officership/team/{officershipteamid}/officership/{officershipid} [delete]
func (s *Store) OfficershipTeamRemoveOfficership(c echo.Context) error {
	teamID, err := strconv.Atoi(c.Param("officershipteamid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get teamid for officership team remove officership: %w", err))
	}

	_, err = s.people.GetOfficershipTeam(c.Request().Context(), people.OfficershipTeamGetDTO{TeamID: teamID})
	if err != nil {
		return fmt.Errorf("failed to get team for officership team remove officership: %w", err)
	}

	officershipID, err := strconv.Atoi(c.Param("officershipid"))
	if err != nil {
		return fmt.Errorf("failed to get officershipid for officership team remove officership: %w", err)
	}

	_, err = s.people.GetOfficership(c.Request().Context(), people.OfficershipGetDTO{OfficershipID: officershipID})
	if err != nil {
		return fmt.Errorf("failed to get officership for officership team remove officership: %w", err)
	}

	officershipTeamMember := people.OfficershipTeamMemberGetDeleteDTO{
		TeamID:    teamID,
		OfficerID: officershipID,
	}

	_, err = s.people.GetOfficershipTeamMember(c.Request().Context(), officershipTeamMember)
	if err != nil {
		return fmt.Errorf("failed to get officership team member for officership team remove officership: %w", err)
	}

	err = s.people.DeleteOfficershipTeamMember(c.Request().Context(), officershipTeamMember)
	if err != nil {
		return fmt.Errorf("failed to remove officership team member for officership team remove officership: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}
