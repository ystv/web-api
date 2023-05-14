package people

import (
	"context"
	"fmt"
)

var _ PermissionRepo = &Store{}

func (s *Store) ListAllPermissions(ctx context.Context) ([]Permission, error) {
	var p []Permission
	for _, permission := range p {
		err := s.db.SelectContext(ctx, &permission, `
			SELECT permission_id, name, description
			FROM people.permissions;`)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions: %w", err)
		}
	}
	return p, nil
}

func (s *Store) ListPermissionMembersByID(ctx context.Context, permissionID int) ([]User, error) {
	var u []User
	err := s.db.GetContext(ctx, &u,
		`SELECT u.user_id, username, email, first_name, last_name, nickname, avatar
		FROM people.users u
		INNER JOIN people.role_members rm ON u.user_id = rm.user_id
		INNER JOIN people.role_permissions p ON rm.role_id = p.role_id
		WHERE permission_id = $1;`, permissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meta: %w", err)
	}
	for _, user := range u {
		if user.Avatar != "" {
			// TODO: sort this out
			user.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + user.Avatar
		}
	}
	return u, nil
}
