package people

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

var _ UserRepo = &Store{}

// GetUserFull will return all user information to be used for profile and management.
func (s *Store) GetUserFull(ctx context.Context, userID int) (UserFull, error) {
	var u UserFull

	err := s.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname,
		avatar, use_gravatar, last_login, created_at, created_by, updated_at, updated_by,
		deleted_at, deleted_by
		FROM people.users
		WHERE user_id = $1
		LIMIT 1;`, userID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get user meta: %w", err)
	}

	err = s.db.SelectContext(ctx, &u.Roles,
		`SELECT r.role_id, r.name, r.description
	FROM people.roles r
	INNER JOIN people.role_members rm ON rm.role_id = r.role_id
	WHERE user_id = $1;`, userID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get roles: %w", err)
	}

	for idx := range u.Roles {
		err := s.db.SelectContext(ctx, &u.Roles[idx].Permissions,
			`SELECT p.permission_id, p.name, p.description
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		WHERE rp.role_id = $1;`, u.Roles[idx].RoleID)
		if err != nil {
			return UserFull{}, fmt.Errorf("failed to get permissions from roles: %w", err)
		}
	}

	if u.Avatar != "" {
		// TODO: sort this out
		u.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + u.Avatar
	}
	if u.UseGravatar {
		hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
		u.Avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%s", hex.EncodeToString(hash[:]))
	}

	return u, nil
}

// GetUser returns basic user information to be used for other services.
func (s *Store) GetUser(ctx context.Context, userID int) (User, error) {
	var u User

	err := s.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname, avatar, use_gravatar
		FROM people.users
		WHERE user_id = $1;`, userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user meta: %w", err)
	}

	err = s.db.SelectContext(ctx, &u.Permissions,
		`SELECT p.permission_id, p.name
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1;`, userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get permissions: %w", err)
	}

	if u.Avatar != "" {
		// TODO sort this out
		u.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + u.Avatar
	}

	if u.UseGravatar {
		hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
		u.Avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%s", hex.EncodeToString(hash[:]))
	}

	return u, nil
}

// ListAllUsers returns all users.
// It doesn't return the full User object.
// Returns user_id, avatar, nickname, first_name, last_name.
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (s *Store) ListAllUsers(ctx context.Context) ([]User, error) {
	var u []User

	err := s.db.SelectContext(ctx, &u,
		`SELECT user_id, avatar, use_gravatar, nickname, first_name, last_name
		FROM people.users;`)
	if err != nil {
		return nil, fmt.Errorf("fialed to list all users: %w", err)
	}

	return u, nil
}
