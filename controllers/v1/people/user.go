package people

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/people"
)

// UserByID finds a user by ID
// @Summary Get a user by ID
// @Description Get a basic user object by ID.
// @ID get-user-id
// @Tags people-user
// @Produce json
// @Param userid path int true "User ID"
// @Success 200 {object} people.User
// @Router /v1/internal/people/user/{userid} [get]
func (r *Repo) UserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	p, err := r.people.GetUser(c.Request().Context(), id)
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
func (r *Repo) UserByIDFull(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	p, err := r.people.GetUserFull(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("UserByIDFull failed to get user: %w", err)
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
func (r *Repo) UserByToken(c echo.Context) error {
	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByToken failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p, err := r.people.GetUser(c.Request().Context(), claims.UserID)
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
func (r *Repo) UserByTokenFull(c echo.Context) error {
	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	p, err := r.people.GetUserFull(c.Request().Context(), claims.UserID)
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed getting user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) AddUser(c echo.Context) error {
	u := people.User{}
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
func (r *Repo) ListAllPeople(c echo.Context) error {
	p, err := r.people.ListAllUsers(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllPeople failed to get all users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}
