package people

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type (
	Repo interface {
		UserRepo
		RoleRepo
		PermissionRepo
	}

	// UserRepo defines all user interactions
	UserRepo interface {
		GetUser(ctx context.Context, userID int) (User, error)
		GetUserFull(ctx context.Context, userID int) (UserFull, error)
		GetUserByEmail(ctx context.Context, email string) (User, error)
		GetUserByEmailFull(ctx context.Context, email string) (UserFull, error)
		ListAllUsers(ctx context.Context) ([]User, error)
	}

	// RoleRepo defines all role interaction
	RoleRepo interface {
		ListAllRolesWithPermissions(ctx context.Context) ([]RoleWithPermissions, error)
		ListAllRolesWithCount(ctx context.Context) ([]RoleWithCount, error)
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
		db          *sqlx.DB
		cdn         *s3.S3
		cdnEndpoint string
	}
)

type (
	// User represents a user object to be used when not all data is required
	User struct {
		UserID      int          `db:"user_id" json:"id"`
		Username    string       `db:"username" json:"username,omitempty"`
		Email       string       `db:"email" json:"email,omitempty"`
		Nickname    string       `db:"nickname" json:"nickname"`
		Avatar      string       `db:"avatar" json:"avatar"`
		UseGravatar bool         `db:"use_gravatar" json:"useGravatar"`
		FirstName   string       `db:"first_name" json:"firstName"`
		LastName    string       `db:"last_name" json:"lastName"`
		Permissions []Permission `json:"permissions,omitempty"`
	}

	// UserFull represents a user and all columns
	UserFull struct {
		User
		LastLogin null.Time `db:"last_login" json:"lastLogin,omitempty"`
		Enabled   bool      `db:"enabled" json:"enabled"`
		CreatedAt null.Time `db:"created_at" json:"createdAt,omitempty"`
		CreatedBy int       `db:"created_by" json:"createdBy,omitempty"`
		UpdatedAt null.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy null.Int  `db:"updated_by" json:"updatedBy,omitempty"`
		DeletedAt null.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy null.Int  `db:"deleted_by" json:"deletedBy,omitempty"`
		Roles     []Role    `json:"roles,omitempty"`
	}

	UserFullPagination struct {
		Users     []UserFull `json:"users"`
		FullCount int        `json:"fullCount"`
	}

	Role struct {
		RoleID      int    `db:"role_id" json:"id"`
		Name        string `db:"name" json:"name"`
		Description string `db:"description" json:"description"`
	}

	// RoleWithPermissions represents a "group" of permissions where multiple users
	// can have this role, and they will inherit these permissions.
	RoleWithPermissions struct {
		Role
		Permissions []Permission `json:"permissions"`
	}

	// RoleWithCount represents a "group" of permissions where multiple users
	// can have this role, and they will inherit these permissions.
	RoleWithCount struct {
		Role
		Users       int `db:"users" json:"users"`
		Permissions int `db:"permissions" json:"permissions"`
	}

	RoleFull struct {
		Role
		Permissions []Permission `json:"permissions"`
		Users       []User       `json:"users"`
	}

	// Permission represents an individual permission. Attempting to implement some RBAC here.
	Permission struct {
		PermissionID int    `db:"permission_id" json:"id"`
		Name         string `db:"name" json:"name"`
		Description  string `db:"description" json:"description,omitempty"`
	}
)

// NewStore creates a new store
func NewStore(db *sqlx.DB, cdn *s3.S3, cdnEndpoint string) Repo {
	return &Store{db: db, cdn: cdn, cdnEndpoint: cdnEndpoint}
}
