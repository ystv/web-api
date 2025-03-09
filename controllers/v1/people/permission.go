package people

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

// ListAllPermissions handles listing all permissions
//
// @Summary List all permissions
// @ID get-people-permissions
// @Tags people-permissions
// @Produce json
// @Success 200 {array} people.Permission
// @Router /v1/internal/people/permissions [get]
func (s *Store) ListAllPermissions(c echo.Context) error {
	p, err := s.people.ListAllPermissions(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllPermissions failed to get permissions: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// ListAllPermissionsWithRolesCount handles listing all permissions
//
// @Summary List all permissions with roles count
// @ID get-people-permissions-count
// @Tags people-permissions
// @Produce json
// @Success 200 {array} people.PermissionWithRolesCount
// @Router /v1/internal/people/permissions/count [get]
func (s *Store) ListAllPermissionsWithRolesCount(c echo.Context) error {
	p, err := s.people.ListAllPermissions(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListAllPermissions failed to get permissions: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// ListPermissionMembersByID handles listing all members of a certain permission
//
// @Summary List all users of a given permission
// @ID get-people-permission-members
// @Tags people-permissions
// @Produce json
// @Param permissionid path int true "Permission ID"
// @Success 200 {array} people.Permission
// @Router /v1/internal/people/permission/{permissionid}/members [get]
func (s *Store) ListPermissionMembersByID(c echo.Context) error {
	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}

	p, err := s.people.ListPermissionMembersByID(c.Request().Context(), permissionID)
	if err != nil {
		err = fmt.Errorf("ListPermissionMembersByID failed to get users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// GetPermissionByID handles getting a single permission
//
// @Summary Get a single permission based on the permission id
// @ID get-people-permission
// @Tags people-permissions
// @Produce json
// @Param permissionid path int true "Permission ID"
// @Success 200 {object} people.Permission
// @Router /v1/internal/people/permission/{permissionid} [get]
func (s *Store) GetPermissionByID(c echo.Context) error {
	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}

	p, err := s.people.GetPermission(c.Request().Context(), permissionID)
	if err != nil {
		err = fmt.Errorf("GetPermission failed to get permission: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// GetPermissionByIDWithRolesCount handles getting a single permission with roles count
//
// @Summary Get a single permission based on the permission id with roles count
// @ID get-people-permission-count
// @Tags people-permissions
// @Produce json
// @Param permissionid path int true "Permission ID"
// @Success 200 {object} people.PermissionWithRolesCount
// @Router /v1/internal/people/permission/{permissionid}/count [get]
func (s *Store) GetPermissionByIDWithRolesCount(c echo.Context) error {
	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}

	p, err := s.people.GetPermissionWithRolesCount(c.Request().Context(), permissionID)
	if err != nil {
		err = fmt.Errorf("GetPermission failed to get permission: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// AddPermission handles creating a permission
//
// @Summary Create a permission
// @ID add-people-permission
// @Tags people-permissions
// @Produce json
// @Param permission body people.PermissionAddEditDTO true "Permission object"
// @Success 201 {object} people.Permission
// @Router /v1/internal/people/permission [post]
func (s *Store) AddPermission(c echo.Context) error {
	var p people.PermissionAddEditDTO

	err := c.Bind(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get permission from request in AddPermission: %w", err))
	}

	permissionAdded, err := s.people.AddPermission(c.Request().Context(), p)
	if err != nil {
		err = fmt.Errorf("AddPermission failed to add permission: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, permissionAdded)
}

// EditPermission handles editing a permission
//
// @Summary Edits a permission
// @ID edit-people-permission
// @Tags people-permissions
// @Produce json
// @Param permissionid path int true "Permission ID"
// @Param permission body people.PermissionAddEditDTO true "Permission object"
// @Success 200 {object} people.Permission
// @Router /v1/internal/people/permission/{permissionid} [put]
func (s *Store) EditPermission(c echo.Context) error {
	var p people.PermissionAddEditDTO

	err := c.Bind(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get permission from request in AddPermission: %w", err))
	}

	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}

	permissionAdded, err := s.people.EditPermission(c.Request().Context(), permissionID, p)
	if err != nil {
		err = fmt.Errorf("EditPermission failed to edit permission: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, permissionAdded)
}

// DeletePermission handles deleting a permission
//
// @Summary Deletes a permission and links to roles
// @ID delete-people-permission
// @Tags people-permissions
// @Produce json
// @Param permissionid path int true "Permission ID"
// @Success 204
// @Router /v1/internal/people/permission/{permissionid} [delete]
func (s *Store) DeletePermission(c echo.Context) error {
	permissionID, err := strconv.Atoi(c.Param("permissionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}

	err = s.people.DeletePermission(c.Request().Context(), permissionID)
	if err != nil {
		err = fmt.Errorf("DeletePermission failed to delete permission: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
