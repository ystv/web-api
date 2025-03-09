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
		CountUsersAll(ctx context.Context) (CountUsers, error)
		ListAllUsers(ctx context.Context) ([]User, error)
		GetUser(ctx context.Context, userID int) (User, error)
		GetUserFull(ctx context.Context, userID int) (UserFull, error)
		GetUserByEmail(ctx context.Context, email string) (User, error)
		GetUserByEmailFull(ctx context.Context, email string) (UserFull, error)
		GetUsersPagination(ctx context.Context, size, page int, search, sortBy, direction, enabled,
			deleted string) ([]UserFull, int, error)
	}

	// RoleRepo defines all role interaction
	RoleRepo interface {
		ListAllRolesWithPermissions(ctx context.Context) ([]RoleWithPermissions, error)
		ListAllRolesWithCount(ctx context.Context) ([]RoleWithCount, error)
		GetRoleFull(ctx context.Context, roleID int) (RoleFull, error)
		ListRoleMembersByID(ctx context.Context, roleID int) ([]User, int, error)
		ListRolePermissionsByID(ctx context.Context, roleID int) ([]Permission, error)
	}

	// PermissionRepo defines all permission interactions
	PermissionRepo interface {
		ListAllPermissions(ctx context.Context) ([]Permission, error)
		ListPermissionMembersByID(ctx context.Context, permissionID int) ([]User, error)
		GetPermission(ctx context.Context, permissionID int) (Permission, error)
		GetPermissionWithRolesCount(ctx context.Context, permissionID int) (PermissionWithRolesCount, error)
		AddPermission(ctx context.Context, permission PermissionAddEditDTO) (Permission, error)
		EditPermission(ctx context.Context, permissionID int, permission PermissionAddEditDTO) (Permission, error)
		DeletePermission(ctx context.Context, permissionID int) error
	}

	// Store contains our dependency
	Store struct {
		db          *sqlx.DB
		cdn         *s3.S3
		cdnEndpoint string
	}

	CountUsers struct {
		TotalUsers             int `db:"total_users" json:"totalUsers"`
		ActiveUsers            int `db:"active_users" json:"activeUsers"`
		ActiveUsersPast24Hours int `db:"active_users_past_24_hours" json:"activeUsersPast24Hours"`
		ActiveUsersPastYear    int `db:"active_users_past_year" json:"activeUsersPastYear"`
	}

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

	// Permission represents an individual permission.
	Permission struct {
		PermissionID int    `db:"permission_id" json:"id"`
		Name         string `db:"name" json:"name"`
		Description  string `db:"description" json:"description,omitempty"`
	}

	// PermissionWithRolesCount represents an individual permission with a count of how many roles ues this.
	PermissionWithRolesCount struct {
		PermissionID int    `db:"permission_id" json:"id"`
		Name         string `db:"name" json:"name"`
		Description  string `db:"description" json:"description,omitempty"`
		Roles        int    `db:"roles" json:"roles"`
	}

	// PermissionAddEditDTO represents a permission to be added or edited.
	PermissionAddEditDTO struct {
		Name        string `db:"name" json:"name"`
		Description string `db:"description" json:"description,omitempty"`
	}
)

// NewStore creates a new store
func NewStore(db *sqlx.DB, cdn *s3.S3, cdnEndpoint string) Repo {
	return &Store{db: db, cdn: cdn, cdnEndpoint: cdnEndpoint}
}
