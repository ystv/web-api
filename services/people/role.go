package people

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Store) ListAllRolesWithPermissions(ctx context.Context) ([]RoleWithPermissions, error) {
	var r []RoleWithPermissions
	//nolint:musttag
	err := s.db.SelectContext(ctx, &r, `
		SELECT role_id, name, description
		FROM people.roles;`)
	if err != nil {
		return nil, fmt.Errorf("failed to select roles: %w", err)
	}

	for _, role := range r {
		err = s.db.SelectContext(ctx, &role.Permissions, `
			SELECT perm.permission_id, name, description
			FROM people.permissions perm
			INNER JOIN people.role_permissions role ON perm.permission_id = role.permission_id
			WHERE role_id = $1;`, role.RoleID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions for role \"%d\": %w", role.RoleID, err)
		}
	}

	return r, nil
}

func (s *Store) ListAllRolesWithCount(ctx context.Context) ([]RoleWithCount, error) {
	var r []RoleWithCount
	//nolint:musttag
	err := s.db.SelectContext(ctx, &r, `SELECT r.*, COUNT(DISTINCT rm.user_id) AS users, 
        COUNT(DISTINCT rp.permission_id) AS permissions
		FROM people.roles r
		LEFT JOIN people.role_members rm ON r.role_id = rm.role_id
		LEFT JOIN people.role_permissions rp ON r.role_id = rp.role_id
		GROUP BY r, r.role_id, name, description
		ORDER BY r.name`)
	if err != nil {
		return nil, fmt.Errorf("failed to select roles: %w", err)
	}

	return r, nil
}

// GetRoleFull returns all users and permissions of a certain role.
// It doesn't return the full User object.
// Returns user_id, avatar, nickname, first_name, last_name.
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (s *Store) GetRoleFull(ctx context.Context, roleID int) (RoleFull, error) {
	var r RoleFull
	//nolint:musttag
	err := s.db.GetContext(ctx, &r, `
		SELECT role_id, name, description
		FROM people.roles
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return RoleFull{}, fmt.Errorf("failed to select roles: %w", err)
	}

	err = s.db.SelectContext(ctx, &r.Permissions, `
		SELECT perm.permission_id, name, description
		FROM people.permissions perm
		INNER JOIN people.role_permissions role ON perm.permission_id = role.permission_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return RoleFull{}, fmt.Errorf("failed to get permissions for role \"%d\": %w", r.RoleID, err)
	}

	//nolint:musttag
	err = s.db.SelectContext(ctx, &r.Users, `
		SELECT u.user_id, avatar, nickname, first_name, last_name
		FROM people.users u
		INNER JOIN people.role_members member ON u.user_id = member.user_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return RoleFull{}, fmt.Errorf("failed to get users for role: %w", err)
	}

	return r, nil
}

// ListRoleMembersByID returns all users who have a certain role.
// It doesn't return the full User object.
// Returns user_id, avatar, nickname, first_name, last_name.
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (s *Store) ListRoleMembersByID(ctx context.Context, roleID int) ([]User, int, error) {
	err := s.validateRole(ctx, roleID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to find role: %w", err)
	}

	var u []User
	//nolint:musttag
	err = s.db.SelectContext(ctx, &u,
		`SELECT u.user_id, avatar, nickname, first_name, last_name
		FROM people.users u
		INNER JOIN people.role_members rm ON u.user_id = rm.user_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to list role users: %w", err)
	}

	return u, http.StatusOK, nil
}

func (s *Store) ListRolePermissionsByID(ctx context.Context, roleID int) ([]Permission, error) {
	var p []Permission

	err := s.db.SelectContext(ctx, &p, `
		SELECT perms.permission_id, perms.name, perms.description
		FROM people.permissions perms
		INNER JOIN people.role_permissions rp ON perms.permission_id = rp.permission_id
		WHERE rp.role_id = $1`, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for role \"%d\": %w", roleID, err)
	}

	return p, nil
}

func (s *Store) validateRole(ctx context.Context, roleID int) error {
	var r RoleFull
	//nolint:musttag
	err := s.db.GetContext(ctx, &r, `
		SELECT role_id, name, description
		FROM people.roles
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return fmt.Errorf("failed to get roles: %w", err)
	}

	return nil
}
