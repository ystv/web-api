package people

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils"
)

// ListAllRolesWithPermissions handles listing all roles
//
// @Summary List all roles
// @ID get-people-roles-with-permissions
// @Tags people-role
// @Produce json
// @Success 200 {array} people.RoleWithPermissions
// @Router /v1/internal/people/roles [get]
func (s *Store) ListAllRolesWithPermissions(c echo.Context) error {
	r, err := s.people.ListAllRolesWithPermissions(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllRolesWithPermissions failed to get roles: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(r))
}

// ListAllRolesWithCount handles listing all roles
//
// @Summary List all roles with count
// @ID get-people-roles-count
// @Tags people-role
// @Produce json
// @Success 200 {array} people.RoleWithCount
// @Router /v1/internal/people/roles/count [get]
func (s *Store) ListAllRolesWithCount(c echo.Context) error {
	r, err := s.people.ListAllRolesWithCount(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllRolesWithCount failed to get roles: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(r))
}

// GetRoleFull handles Getting a certain role and all users and permissions
//
// @Summary List all users and permissions of a given role
// @ID get-people-role-full
// @Tags people-role
// @Produce json
// @Param roleid path int true "Role ID"
// @Success 200 {object} people.RoleFull
// @Router /v1/internal/people/role/{roleid}/full [get]
func (s *Store) GetRoleFull(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	r, err := s.people.GetRoleFull(c.Request().Context(), roleID)
	if err != nil {
		err = fmt.Errorf("GetRoleFull failed to get users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, r)
}

// ListRoleMembersByID handles listing all members of a certain role
//
// @Summary List all users of a given role
// @ID get-people-role-members
// @Tags people-role
// @Produce json
// @Param roleid path int true "Role ID"
// @Success 200 {array} people.User
// @Failure 404 {object} utils.HTTPError "Role Not Found"
// @Failure 500 {object} utils.HTTPError "Server Side Role Error"
// @Router /v1/internal/people/role/{roleid}/members [get]
func (s *Store) ListRoleMembersByID(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	u, status, err := s.people.ListRoleMembersByID(c.Request().Context(), roleID)
	if err != nil {
		err = fmt.Errorf("ListRoleMembersByID failed to get users: %w", err)
		return echo.NewHTTPError(status, err)
	}

	return c.JSON(http.StatusOK, u)
}

// ListRolePermissionsByID handles listing permissions of a certain role
//
// @Summary List permissions of a given role
// @ID get-people-role-permissions
// @Tags people-role
// @Produce json
// @Param roleid path int true "Role ID"
// @Success 200 {array} people.Permission
// @Router /v1/internal/people/role/{roleid}/permissions [get]
func (s *Store) ListRolePermissionsByID(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	p, err := s.people.ListRolePermissionsByID(c.Request().Context(), roleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("ListRolePermissionsByID failed to get permission: %w", err))
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}
