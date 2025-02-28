package people

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

// UserStats lists the user stats
// @Summary Get user stats
// @Description Get an overview of users.
// @ID get-user-stats
// @Tags people-user
// @Produce json
// @Success 200 {object} people.CountUsers
// @Router /v1/internal/people/users/stats [get]
func (s *Store) UserStats(c echo.Context) error {
	stat, err := s.people.CountUsersAll(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("UserStats failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, stat)
}

// UserByID finds a user by ID
// @Summary Get a user by ID
// @Description Get a basic user object by ID.
// @ID get-user-id
// @Tags people-user
// @Produce json
// @Param userid path int true "User ID"
// @Success 200 {object} people.User
// @Router /v1/internal/people/user/{userid} [get]
func (s *Store) UserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	p, err := s.people.GetUser(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("UserByID failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// UserByIDFull finds a user by ID returning all info
// @Summary Get a full user by ID
// @Description Get a complete user object by ID.
// @ID get-user-id-full
// @Tags people-user
// @Produce json
// @Param userid path int true "User ID"
// @Success 200 {object} people.User
// @Router /v1/internal/people/user/{userid}/full [get]
func (s *Store) UserByIDFull(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	p, err := s.people.GetUserFull(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("UserByIDFull failed to get user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// UserByEmail finds a user by email
// @Summary Get a user by email
// @Description Get a basic user object by email.
// @ID get-user-email
// @Tags people-user
// @Produce json
// @Param email path int true "Email"
// @Success 200 {object} people.User
// @Router /v1/internal/people/user/{email} [get]
func (s *Store) UserByEmail(c echo.Context) error {
	temp := c.Param("email")
	email, err := url.QueryUnescape(temp)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	p, err := s.people.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		err = fmt.Errorf("UserByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// UserByEmailFull finds a user by email returning all info
// @Summary Get a full user by email
// @Description Get a complete user object by email.
// @ID get-user-email-full
// @Tags people-user
// @Produce json
// @Param email path int true "Email"
// @Success 200 {object} people.User
// @Router /v1/internal/people/user/{email}/full [get]
func (s *Store) UserByEmailFull(c echo.Context) error {
	temp := c.Param("email")
	email, err := url.QueryUnescape(temp)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	p, err := s.people.GetUserByEmailFull(c.Request().Context(), email)
	if err != nil {
		err = fmt.Errorf("UserByEmailFull failed to get user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// UserByToken finds a user by their JWT token
// @Summary Get a user by token
// @Description Get a basic user object by JWT token generated by web-auth.
// @ID get-user-token
// @Tags people-user
// @Produce json
// @Success 200 {object} people.User
// @Router /v1/internal/people/user [get]
func (s *Store) UserByToken(c echo.Context) error {
	claims, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByToken failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	p, err := s.people.GetUser(c.Request().Context(), claims.UserID)
	if err != nil {
		err = fmt.Errorf("UserByToken failed getting user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// UserByTokenFull finds a user by their JWT token returning all info
// @Summary Get a full user by token
// @Description Get a complete user object by JWT token generated by web-auth.
// @ID get-user-token-full
// @Tags people-user
// @Produce json
// @Success 200 {object} people.UserFull
// @Router /v1/internal/people/user/full [get]
func (s *Store) UserByTokenFull(c echo.Context) error {
	claims, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	p, err := s.people.GetUserFull(c.Request().Context(), claims.UserID)
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed getting user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

func (s *Store) AddUser(c echo.Context) error {
	var u people.User

	err := c.Bind(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get user from request in AddUser: %w", err))
	}

	return c.JSON(http.StatusOK, u)
}

// ListAllPeople handles listing all users
//
// @Summary List all users
// @ID get-people-users-all
// @Tags people-users
// @Produce json
// @Success 200 {array} people.User
// @Router /v1/internal/people/users [get]
func (s *Store) ListAllPeople(c echo.Context) error {
	p, err := s.people.ListAllUsers(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllPeople failed to get all users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}
