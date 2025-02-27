package people

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils"
)

// ListAllRoles handles listing all roles
//
// @Summary List all roles
// @ID get-people-roles
// @Tags people-roles
// @Produce json
// @Success 200 {array} people.Role
// @Router /v1/internal/people/role [get]
func (r *Repo) ListAllRoles(c echo.Context) error {
	p, err := r.people.ListAllRolesWithPermissions(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllRolesWithPermissions failed to get roles: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// ListRoleMembersByID handles listing all members of a certain role
//
// @Summary List all users of a given role
// @ID get-people-role-members
// @Tags people-roles
// @Produce json
// @Param roleId path int true "Role ID"
// @Success 200 {array} people.Role
// @Router /v1/internal/people/role/{roleId}/members [get]
func (r *Repo) ListRoleMembersByID(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	p, err := r.people.ListRoleMembersByID(c.Request().Context(), roleID)
	if err != nil {
		err = fmt.Errorf("ListRoleMembersByID failed to get users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// ListRolePermissionsByID handles listing all permissions of a certain role
//
// @Summary List all permissions of a given role
// @ID get-people-role-permissions
// @Tags people-roles
// @Produce json
// @Param roleId path int true "Role ID"
// @Success 200 {array} people.Role
// @Router /v1/internal/people/role/{roleId}/permissions [get]
func (r *Repo) ListRolePermissionsByID(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	p, err := r.people.ListRolePermissionsByID(c.Request().Context(), roleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("ListRolePermissionsByID failed to get permission: %w", err))
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}
