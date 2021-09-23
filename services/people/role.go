package people

import (
	"context"
	"fmt"
)

var _ RoleRepo = &Store{}

func (s *Store) ListAllRoles(ctx context.Context) ([]Role, error) {
	r := []Role{}
	err := s.db.SelectContext(ctx, &r, `
		SELECT name, description
		FROM people.roles;`)
	if err != nil {
		return nil, fmt.Errorf("failed to select roles: %w", err)
	}
	for _, role := range r {
		err = s.db.SelectContext(ctx, &role, `
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

// ListRole returns all users who have a certain role
// It doesn't return the full User object
// Returns user_id, avatar, nickname, first_name, last_name
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (m *Store) ListRoleMembersByID(ctx context.Context, roleID int) ([]User, error) {
	u := []User{}
	err := m.db.SelectContext(ctx, &u,
		`SELECT u.user_id, avatar, nickname, first_name, last_name
		FROM people.users u
		INNER JOIN people.role_members rm ON u.user_id = rm.user_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list role users: %w", err)
	}
	return u, nil
}

func (m *Store) ListRolePermissionsByID(ctx context.Context, roleID int) ([]Permission, error) {
	p := []Permission{}
	for _, permission := range p {
		err := m.db.SelectContext(ctx, &permission, `
			SELECT perm.permission_id, name, description
			FROM people.permissions perm
			INNER JOIN people.role_permissions role ON perm.permission_id = role.permission_id
			WHERE role_id = $1;`, roleID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions for role \"%d\": %w", roleID, err)
		}
	}
	return p, nil
}
