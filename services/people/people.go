package people

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// UserRepo defines all user interactions
	UserRepo interface {
		GetUser(ctx context.Context, userID int) (User, error)
		GetUserFull(ctx context.Context, userID int) (UserFull, error)
		ListAllUsers(ctx context.Context) ([]User, error)
	}

	// RoleRepo defines all role interaction
	RoleRepo interface {
		ListAllRoles(ctx context.Context) ([]Role, error)
		ListRoleMembersByID(ctx context.Context, roleID int) ([]User, error)
		ListRolePermissionsByID(ctx context.Context, roleID int) ([]Permission, error)
	}

	// PermissionRepo defines all permission interactions
	PermissionRepo interface {
		ListAllPermissions(ctx context.Context) ([]Permission, error)
		ListPermissionMembersByID(ctx context.Context, permissionID int) ([]User, error)
	}

	// Store contains our dependency
	Store struct {
		db *sqlx.DB
	}
)

//TODO Sort out pointers. They are currently here so when the json is being marshalled it will "omitempty"

type (
	//User represents a user object to be used when not all data is required
	User struct {
		UserID      int          `db:"user_id" json:"id"`
		Username    string       `db:"username" json:"username,omitempty"`
		Email       string       `db:"email" json:"email,omitempty"`
		Nickname    string       `db:"nickname" json:"nickname"`
		Avatar      string       `db:"avatar" json:"avatar"`
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
		UpdatedBy *int       `db:"updated_by" json:"updatedBy,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy *int       `db:"deleted_by" json:"deletedBy,omitempty"`
		Roles     []Role     `json:"roles,omitempty"`
	}
	// Role represents a "group" of permissions where multiple users
	// can have this role and they will inherit these permissions.
	Role struct {
		RoleID      int          `db:"role_id" json:"id"`
		Name        string       `db:"name" json:"name"`
		Description string       `db:"description" json:"description"`
		Permissions []Permission `json:"permissions"`
	}

	// Permission represents an individual permission. Attempting to implement some RBAC here.
	Permission struct {
		PermissionID int    `db:"permission_id" json:"id"`
		Name         string `db:"name" json:"name"`
		Description  string `db:"description" json:"description,omitempty"`
	}
)

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}
