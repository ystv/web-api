package people

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

//TODO Sort out pointers. They are currently here so when the json is being marshalled it will "omitempty"

type (
	//User represents a user object to be used when not all data is required
	User struct {
		ID          int          `db:"user_id" json:"id"`
		Username    string       `db:"username" json:"username,omitempty"`
		Email       string       `db:"email" json:"email,omitempty"`
		Nickname    string       `db:"nickname" json:"nickname"`
		Avatar      null.String  `db:"avatar" json:"avatar"`
		FirstName   string       `db:"first_name" json:"firstName"`
		LastName    string       `db:"last_name" json:"lastName"`
		Permissions []Permission `json:"permissions,omitempty"`
	}
	// UserFull represents a user and all columns
	UserFull struct {
		User
		LastLogin *time.Time `db:"last_login" json:"lastLogin,omitempty"`
		CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
		CreatedBy int        `db:"created_by" json:"createdBy,omitempty"`
		UpdatedAt *null.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy *null.Int  `db:"updated_by" json:"updatedBy,omitempty"`
		DeletedAt *null.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy *null.Int  `db:"deleted_by" json:"deletedBy,omitempty"`
		Roles     []Role     `json:"roles,omitempty"`
	}
	// Role represents a "group" of permissions where multiple users
	// can have this role and they will inherit these permissions.
	Role struct {
		ID          int          `db:"role_id" json:"id"`
		Name        string       `db:"name" json:"name"`
		Description null.String  `db:"description" json:"description"`
		Permissions []Permission `json:"permissions"`
	}

	// Permission represents an individual permission. Attempting to implement some RBAC here.
	Permission struct {
		ID          int          `db:"permission_id" json:"id"`
		Name        string       `db:"name" json:"name"`
		Description *null.String `db:"description" json:"description,omitempty"`
	}
)

// Here for validation to ensure we are meeting the interface
var _ UserRepo = &Store{}

// GetFull will return all user information to be used for profile and management.
func (m *Store) GetFull(ctx context.Context, userID int) (*UserFull, error) {
	u := UserFull{}
	err := m.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname,
		avatar, last_login, created_at, created_by, updated_at, updated_by,
		deleted_at, deleted_by
		FROM people.users
		WHERE user_id = $1
		LIMIT 1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meta: %w", err)
	}
	err = m.db.SelectContext(ctx, &u.Roles,
		`SELECT r.role_id, r.name, r.description
	FROM people.roles r
	INNER JOIN people.role_members rm ON rm.role_id = r.role_id
	WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	for idx := range u.Roles {
		err := m.db.SelectContext(ctx, &u.Roles[idx].Permissions,
			`SELECT p.permission_id, p.name, p.description
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		WHERE rp.role_id = $1;`, u.Roles[idx].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions from roles: %w", err)
		}
	}
	if u.Avatar.Valid {
		// TODO sort this out
		u.Avatar = null.StringFrom("https://ystv.co.uk/static/images/members/thumb/" + u.Avatar.String)
	}
	return &u, nil
}

// Get returns basic user information to be used for other services.
func (m *Store) Get(ctx context.Context, userID int) (*User, error) {
	u := User{}
	err := m.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname, avatar
		FROM people.users
		WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meta: %w", err)
	}
	err = m.db.SelectContext(ctx, &u.Permissions,
		`SELECT p.permission_id, p.name
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	if u.Avatar.Valid {
		// TODO sort this out
		u.Avatar = null.StringFrom("https://ystv.co.uk/static/images/members/thumb/" + u.Avatar.String)
	}
	return &u, nil
}
