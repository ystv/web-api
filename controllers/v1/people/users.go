package people

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ListAllPeople handles listing all users
func (r *Repo) ListAllPeople(c echo.Context) error {
	p, err := r.user.ListAll(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllPeople failed to get all users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// ListRoleMembers handles listing all members of a certain role
func (r *Repo) ListRoleMembers(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	p, err := r.user.ListRole(c.Request().Context(), roleID)
	if err != nil {
		err = fmt.Errorf("ListRoleMembers failed to get users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}
