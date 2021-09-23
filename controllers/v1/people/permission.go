package people

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ListAllPermissions handles listing all permissions
//
// @Summary List all permissions
// @ID get-people-permissions
// @Tags people-permissions
// @Produce json
// @Success 200 {array} people.Permission
// @Router /v1/internal/people/permission [get]
func (r *Repo) ListAllPermissions(c echo.Context) error {
	p, err := r.people.ListAllPermissions(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllPermissions failed to get permissions: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// ListPermissionMembersByID handles listing all members of a certain permission
//
// @Summary List all users of a given permission
// @ID get-people-permission-members
// @Tags people-permissions
// @Produce json
// @Param permId path int true "Permission ID"
// @Success 200 {array} people.Permission
// @Router /v1/internal/people/permission/{permId}/members [get]
func (r *Repo) ListPermissionMembersByID(c echo.Context) error {
	permissionID, err := strconv.Atoi(c.Param("permId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}
	p, err := r.people.ListPermissionMembersByID(c.Request().Context(), permissionID)
	if err != nil {
		err = fmt.Errorf("ListPermissionMembersByID failed to get users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}
