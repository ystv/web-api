package people

import (
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
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy int        `db:"updated_by" json:"updatedBy,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy int        `db:"deleted_by" json:"deletedBy,omitempty"`
		Roles     []Role     `json:"roles,omitempty"`
	}
	// Role represents a "group" of permissions where multiple users
	// can have this role and they will inherit these permissions.
	Role struct {
		ID          int          `db:"role_id" json:"roleID"`
		Name        string       `db:"name" json:"name"`
		Description null.String  `db:"description" json:"description"`
		Permissions []Permission `json:"permissions"`
	}

	// Permission represents an individual permission. Attempting to implement some RBAC here.
	Permission struct {
		ID          int          `db:"permission_id" json:"permissionID"`
		Name        string       `db:"name" json:"name"`
		Description *null.String `db:"description" json:"description,omitempty"`
	}
)

// GetFull will return all user information to be used for profile and management.
func GetFull(id int) (*UserFull, error) {
	create := time.Date(2015, 9, 18, 0, 0, 0, 0, time.UTC)
	login := time.Now()
	return &UserFull{
		User: User{
			ID:        3348,
			FirstName: "Rhys",
			LastName:  "Milling",
			Username:  "rhys",
			Nickname:  "Mr. Cool",
			Avatar:    null.StringFrom("https://ystv.co.uk/static/images/members/thumb/3348.jpg"),
			Email:     "rdjm502@york.ac.uk",
		},
		LastLogin: &login,
		CreatedAt: &create,
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

// Get returns basic user information to be used for other services.
func Get(id int) (*User, error) {
	return &User{
		ID:        3348,
		FirstName: "Rhys",
		LastName:  "Milling",
		Nickname:  "Mr. Cool",
		Avatar:    null.StringFrom("https://ystv.co.uk/static/images/members/thumb/3348.jpg"),
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
