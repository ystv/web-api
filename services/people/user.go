package people

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

//TODO Sort out pointers. They are currently here so when the json is being marshalled it will "omitempty"

// User represents a normal user
type User struct {
	ID          int          `db:"user_id" json:"id"`
	Username    string       `db:"username" json:"username,omitempty"`
	Email       string       `db:"email" json:"email,omitempty"`
	DisplayName string       `db:"display_name" json:"displayName"`
	Avatar      null.String  `db:"avatar" json:"avatar"`
	FirstName   string       `db:"first_name" json:"firstName"`
	LastName    string       `db:"last_name" json:"lastName"`
	LastLogin   *time.Time   `db:"last_login" json:"lastLogin,omitempty"`
	CreatedAt   *time.Time   `db:"created_at" json:"createdAt,omitempty"`
	CreatedBy   int          `db:"created_by" json:"createdBy,omitempty"`
	Roles       []Role       `json:"roles,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// Role represents a "group" of permissions where multiple users
// can have this role and they will inherit these permissions.
type Role struct {
	ID          int          `db:"role_id" json:"roleID"`
	Name        string       `db:"name" json:"name"`
	Description null.String  `db:"description" json:"description"`
	Permissions []Permission `json:"permissions"`
}

// Permission represents an individual permission. Attempting to implement some RBAC here.
type Permission struct {
	ID          int          `db:"permission_id" json:"permissionID"`
	Name        string       `db:"name" json:"name"`
	Description *null.String `db:"description" json:"description,omitempty"`
}

// GetFull will return minimal user information to be used for other services.
func GetFull(id int) (*User, error) {
	create := time.Date(2015, 9, 18, 0, 0, 0, 0, time.UTC)
	login := time.Now()
	return &User{
		ID:          1,
		FirstName:   "Rhys",
		LastName:    "Milling",
		Username:    "rhys",
		DisplayName: "Mr. Cool",
		Email:       "rdjm502@york.ac.uk",
		LastLogin:   &login,
		CreatedAt:   &create,
		Roles: []Role{
			{
				ID:   1,
				Name: "Comp team",
				Permissions: []Permission{
					{
						ID:   43,
						Name: "creatorAdmin",
					},
					{
						ID:   67,
						Name: "quoteAdmin",
					},
				},
			},
		},
	}, nil
}

func Get(id int) (*User, error) {
	return &User{
		ID:          1,
		FirstName:   "Rhys",
		LastName:    "Milling",
		DisplayName: "Mr. Cool",
		Permissions: []Permission{
			{
				ID:   45,
				Name: "creatorAdmin",
			},
			{
				ID:   78,
				Name: "emailAccess",
			},
		},
	}, nil
}
