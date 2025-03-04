package people

import (
	"errors"
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
	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByToken failed to get token: %w", err)
		return echo.NewHTTPError(status, err)
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
	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed to get token: %w", err)
		return echo.NewHTTPError(status, err)
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

	return c.JSON(http.StatusNotImplemented, u)
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

// ListPeoplePagination handles listing users with pagination
//
// @Summary List users with pagination
// @ID get-people-users-pagination
// @Tags people-users
// @Produce json
// @Param size path int false "Page size"
// @Param page path int false "Page number"
// @Param search path string false "Search string"
// @Param column path string false "Ordering column"
// @Param direction path string false "Ordering direction"
// @Param enabled path string false "Is user enabled"
// @Param deleted path string false "Is user deleted"
// @Success 200 {array} people.UserFull
// @Router /v1/internal/people/users/pagination [get]
func (s *Store) ListPeoplePagination(c echo.Context) error {
	column := c.QueryParam("column")
	direction := c.QueryParam("direction")
	search := c.QueryParam("search")

	search, err := url.QueryUnescape(search)
	if err != nil {
		return fmt.Errorf("ListPeoplePagination failed to unescape query: %w", err)
	}

	enabled := c.QueryParam("enabled")
	deleted := c.QueryParam("deleted")

	var size, page int

	sizeRaw := c.QueryParam("size")

	if sizeRaw == "all" {
		size = 0
	} else if len(sizeRaw) != 0 {
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest,
				fmt.Errorf("ListPeoplePagination unable to parse page for users: %w", err))
		}

		size, err = strconv.Atoi(sizeRaw)
		//nolint:gocritic
		if err != nil {
			size = 0
		} else if size <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest,
				errors.New("ListPeoplePagination invalid size, must be positive"))
		} else if size != 5 && size != 10 && size != 25 && size != 50 && size != 75 && size != 100 {
			size = 0
		}
	}

	switch column {
	case "userId", "name", "username", "email", "lastLogin":
		switch direction {
		case "asc", "desc":
			break
		default:
			column = ""
			direction = ""
		}
	default:
		column = ""
		direction = ""
	}

	dbUsers, fullCount, err := s.people.GetUsersPagination(c.Request().Context(), size, page, search, column,
		direction, enabled, deleted)
	if err != nil {
		err = fmt.Errorf("ListPeoplePagination failed to get paginated users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	u := people.UserFullPagination{
		Users:     dbUsers,
		FullCount: fullCount,
	}

	return c.JSON(http.StatusOK, u)
}
