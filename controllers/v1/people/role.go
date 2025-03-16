package people

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/people"
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

// AddRole handles creating a role
//
// @Summary Create a role
// @ID add-people-role
// @Tags people-role
// @Produce json
// @Param role body people.RoleAddEditDTO true "Role object"
// @Success 201 {object} people.Role
// @Router /v1/internal/people/role [post]
func (s *Store) AddRole(c echo.Context) error {
	var roleAdd people.RoleAddEditDTO

	err := c.Bind(&roleAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to parse body")
	}

	r1, err := s.people.GetRole(c.Request().Context(), people.RoleGetDTO{Name: roleAdd.Name})
	if err == nil && r1.RoleID > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "role with name \""+roleAdd.Name+"\" already exists")
	}

	role, err := s.people.AddRole(c.Request().Context(), roleAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to add role: %w", err))
	}

	return c.JSON(http.StatusOK, role)
}

// EditRole handles editing a role
//
// @Summary Edit a role
// @ID edit-people-role
// @Tags people-role
// @Param roleid path int true "Role ID"
// @Produce json
// @Param role body people.RoleAddEditDTO true "Role object"
// @Success 200 {object} people.Role
// @Router /v1/internal/people/role/{roleid} [put]
func (s *Store) EditRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return fmt.Errorf("failed to get roleid for edit role: %w", err)
	}

	_, err = s.people.GetRole(c.Request().Context(),
		people.RoleGetDTO{RoleID: roleID})
	if err != nil {
		return fmt.Errorf("failed to get role for edit role: %w", err)
	}

	var roleEdit people.RoleAddEditDTO

	err = c.Bind(&roleEdit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to parse body")
	}

	r1, err := s.people.GetRole(c.Request().Context(), people.RoleGetDTO{Name: roleEdit.Name})
	if err == nil && r1.RoleID > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "role with name \""+roleEdit.Name+"\" already exists")
	}

	role, err := s.people.EditRole(c.Request().Context(), roleID, roleEdit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to edit role: %w", err))
	}

	return c.JSON(http.StatusOK, role)
}

// DeleteRole handles deleting a role
//
// @Summary Deletes role
// @ID delete-people-role
// @Tags people-role
// @Param roleid path int true "role id"
// @Produce json
// @Success 204
// @Router /v1/internal/people/role/{roleid} [delete]
func (s *Store) DeleteRole(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get roleid for role: %w", err))
	}

	_, err = s.people.GetRole(c.Request().Context(), people.RoleGetDTO{RoleID: roleID})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get role for delete role: %w", err))
	}

	err = s.people.RemoveRoleForPermissions(c.Request().Context(), roleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete role permission for delete role: %w", err))
	}

	err = s.people.RemoveRoleForUsers(c.Request().Context(), roleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete role user for delete role: %w", err))
	}

	err = s.people.DeleteRole(c.Request().Context(), roleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete role for delete role: %w", err))
	}

	return c.NoContent(http.StatusNoContent)
}

// RoleAddPermissionFunc handles adding a permission to a role
//
// @Summary Adds a permission to a role
// @ID add-people-role-permission
// @Tags people-role
// @Produce json
// @Param roleid path int true "role id"
// @Param permissionid path int true "permission id"
// @Success 201 {object} people.RolePermission
// @Router /v1/internal/people/role/{roleid}/permission/{permissionid} [post]
func (s *Store) RoleAddPermissionFunc(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get role for roleAddPermission: %w", err))
	}

	_, err = s.people.GetRole(c.Request().Context(), people.RoleGetDTO{RoleID: roleID})
	if err != nil {
		return fmt.Errorf("failed to get role for roleAddPermission: %w", err)
	}

	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return fmt.Errorf("failed to get permissionid for roleAddPermission: %w", err)
	}

	_, err = s.people.GetPermission(c.Request().Context(), permissionID)
	if err != nil {
		return fmt.Errorf("failed to get permission for roleAddPermission: %w", err)
	}

	rolePermission := people.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	_, err = s.people.GetRolePermission(c.Request().Context(), rolePermission)
	if err == nil {
		return errors.New("failed to add rolePermission for roleAddPermission: row already exists")
	}

	rp, err := s.people.AddRolePermission(c.Request().Context(), rolePermission)
	if err != nil {
		return fmt.Errorf("failed to add rolePermission for roleAddPermission: %w", err)
	}

	return c.JSON(http.StatusOK, rp)
}

// RoleRemovePermissionFunc handles removing a permission from a role
//
// @Summary Removes a permission from a role
// @ID remove-people-role-permission
// @Tags people-role
// @Produce json
// @Param roleid path int true "role id"
// @Param permissionid path int true "permission id"
// @Success 204
// @Router /v1/internal/people/role/{roleid}/permission/{permissionid} [delete]
func (s *Store) RoleRemovePermissionFunc(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get roleid for roleRemovePermission: %w", err))
	}

	_, err = s.people.GetRole(c.Request().Context(), people.RoleGetDTO{RoleID: roleID})
	if err != nil {
		return fmt.Errorf("failed to get role for roleRemovePermission: %w", err)
	}

	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return fmt.Errorf("failed to get permissionid for roleRemovePermission: %w", err)
	}

	_, err = s.people.GetPermission(c.Request().Context(), permissionID)
	if err != nil {
		return fmt.Errorf("failed to get permission for roleRemovePermission: %w", err)
	}

	rolePermission := people.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	_, err = s.people.GetRolePermission(c.Request().Context(), rolePermission)
	if err != nil {
		return fmt.Errorf("failed to get rolePermisison for roleRemovePermission: %w", err)
	}

	err = s.people.RemoveRolePermission(c.Request().Context(), rolePermission)
	if err != nil {
		return fmt.Errorf("failed to remove rolePermission for roleRemoveRole: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RoleAddUserFunc handles adding a user to a role
//
// @Summary Adds a user to a role
// @ID add-people-role-user
// @Tags people-role
// @Produce json
// @Param roleid path int true "role id"
// @Param userid path int true "user id"
// @Success 201 {object} people.RoleUser
// @Router /v1/internal/people/role/{roleid}/user/{userid} [post]
func (s *Store) RoleAddUserFunc(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get role for roleAddUser: %w", err))
	}

	_, err = s.people.GetRole(c.Request().Context(), people.RoleGetDTO{RoleID: roleID})
	if err != nil {
		return fmt.Errorf("failed to get user for roleAddUser: %w", err)
	}

	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		return fmt.Errorf("failed to get userID for roleAddUser: %w", err)
	}

	_, err = s.people.GetUser(c.Request().Context(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user for roleAddUser: %w", err)
	}

	roleUser := people.RoleUser{
		RoleID: roleID,
		UserID: userID,
	}

	_, err = s.people.GetRoleUser(c.Request().Context(), roleUser)
	if err == nil {
		return errors.New("failed to add roleUser for roleAddUser: row already exists")
	}

	ru, err := s.people.AddRoleUser(c.Request().Context(), roleUser)
	if err != nil {
		return fmt.Errorf("failed to add roleUser for roleAddUser: %w", err)
	}

	return c.JSON(http.StatusOK, ru)
}

// RoleRemoveUserFunc handles removing a user from a role
//
// @Summary Removes a user from a role
// @ID remove-people-role-user
// @Tags people-role
// @Produce json
// @Param roleid path int true "role id"
// @Param userid path int true "user id"
// @Success 204
// @Router /v1/internal/people/role/{roleid}/user/{userid} [delete]
func (s *Store) RoleRemoveUserFunc(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("roleid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get roleid for roleRemoveUser: %w", err))
	}

	_, err = s.people.GetRole(c.Request().Context(), people.RoleGetDTO{RoleID: roleID})
	if err != nil {
		return fmt.Errorf("failed to get role for roleRemoveUser: %w", err)
	}

	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		return fmt.Errorf("failed to get userID for roleRemoveUser: %w", err)
	}

	_, err = s.people.GetUser(c.Request().Context(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user for roleRemoveUser: %w", err)
	}

	roleUser := people.RoleUser{
		RoleID: roleID,
		UserID: userID,
	}

	_, err = s.people.GetRoleUser(c.Request().Context(), roleUser)
	if err != nil {
		return fmt.Errorf("failed to get roleUser for roleRemoveUser: %w", err)
	}

	err = s.people.RemoveRoleUser(c.Request().Context(), roleUser)
	if err != nil {
		return fmt.Errorf("failed to remove roleUser for roleRemoveUser: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}
