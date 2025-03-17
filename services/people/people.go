package people

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type (
	Repo interface {
		UserRepo
		RoleRepo
		PermissionRepo
		OfficershipRepo
	}

	OfficershipRepo interface {
		OfficershipDBToOfficership(officershipDB OfficershipDB) Officership
		OfficershipMemberDBToOfficershipMember(officershipDB OfficershipMemberDB) OfficershipMember
		CountOfficerships(ctx context.Context) (CountOfficerships, error)
		GetOfficerships(context.Context, OfficershipsStatus) ([]OfficershipDB, error)
		GetOfficership(context.Context, OfficershipGetDTO) (OfficershipDB, error)
		AddOfficership(context.Context, OfficershipAddEditDTO) (OfficershipDB, error)
		EditOfficership(context.Context, int, OfficershipAddEditDTO) (OfficershipDB, error)
		DeleteOfficership(ctx context.Context, officershipID int) error
		GetOfficershipTeams(ctx context.Context) ([]OfficershipTeam, error)
		GetOfficershipTeam(context.Context, OfficershipTeamGetDTO) (OfficershipTeam, error)
		AddOfficershipTeam(context.Context, OfficershipTeamAddEditDTO) (OfficershipTeam, error)
		EditOfficershipTeam(context.Context, int, OfficershipTeamAddEditDTO) (OfficershipTeam, error)
		DeleteOfficershipTeam(ctx context.Context, officershipTeamID int) error
		GetOfficershipTeamMembers(ctx context.Context, officershipTeamID *int,
			officershipStatus OfficershipsStatus) ([]OfficershipTeamMember, error)
		GetOfficershipsNotInTeam(ctx context.Context, officershipTeamID int) ([]OfficershipDB, error)
		GetOfficershipTeamMember(ctx context.Context, officershipTeamMemberGet OfficershipTeamMemberGetDeleteDTO) (OfficershipTeamMember, error)
		AddOfficershipTeamMember(context.Context, OfficershipTeamMemberAddDTO) (OfficershipTeamMember, error)
		DeleteOfficershipTeamMember(context.Context, OfficershipTeamMemberGetDeleteDTO) error
		RemoveTeamForOfficershipTeamMembers(ctx context.Context, officershipTeamID int) error
		GetOfficershipMembers(ctx context.Context, officershipGet *OfficershipGetDTO, userID *int,
			officershipStatus OfficershipsStatus, officershipMemberStatus OfficershipsStatus,
			orderByOfficerName bool) ([]OfficershipMemberDB, error)
		GetOfficershipMember(ctx context.Context, officershipMemberID int) (OfficershipMemberDB, error)
		AddOfficershipMember(context.Context, OfficershipMemberAddEditDTO) (OfficershipMemberDB, error)
		EditOfficershipMember(context.Context, int, OfficershipMemberAddEditDTO) (OfficershipMemberDB, error)
		DeleteOfficershipMember(ctx context.Context, officershipMemberID int) error
		RemoveOfficershipForOfficershipMembers(ctx context.Context, officershipID int) error
		RemoveUserForOfficershipMembers(ctx context.Context, userID int) error
	}

	// UserRepo defines all user interactions
	UserRepo interface {
		CountUsersAll(ctx context.Context) (CountUsers, error)
		ListAllUsers(ctx context.Context) ([]User, error)
		GetUser(ctx context.Context, userID int) (User, error)
		GetUserFull(ctx context.Context, userID int) (UserFullDB, error)
		GetUserByEmail(ctx context.Context, email string) (User, error)
		GetUserByEmailFull(ctx context.Context, email string) (UserFullDB, error)
		GetUsersPagination(ctx context.Context, size, page int, search, sortBy, direction, enabled,
			deleted string) ([]UserFullDB, int, error)
	}

	// RoleRepo defines all role interaction
	RoleRepo interface {
		ListAllRolesWithPermissions(ctx context.Context) ([]RoleWithPermissions, error)
		ListAllRolesWithCount(ctx context.Context) ([]RoleWithCount, error)
		GetRole(ctx context.Context, roleGetDTO RoleGetDTO) (Role, error)
		GetRoleFull(ctx context.Context, roleID int) (RoleFull, error)
		ListRoleMembersByID(ctx context.Context, roleID int) ([]User, int, error)
		ListRolePermissionsByID(ctx context.Context, roleID int) ([]Permission, error)
		AddRole(ctx context.Context, roleAdd RoleAddEditDTO) (Role, error)
		EditRole(ctx context.Context, roleID int, roleEdit RoleAddEditDTO) (Role, error)
		DeleteRole(ctx context.Context, roleID int) error
		RemoveRoleForPermissions(ctx context.Context, roleID int) error
		RemoveRoleForUsers(ctx context.Context, roleID int) error
		GetRoleUser(context.Context, RoleUser) (RoleUser, error)
		GetUsersNotInRole(ctx context.Context, roleID int) ([]User, error)
		AddRoleUser(context.Context, RoleUser) (RoleUser, error)
		RemoveRoleUser(context.Context, RoleUser) error
		RemoveUserForRoles(context.Context, User) error
		GetPermissionsForRole(context.Context, int) ([]Permission, error)
		GetRolesForPermission(context.Context, int) ([]Role, error)
		GetRolePermission(context.Context, RolePermission) (RolePermission, error)
		GetPermissionsNotInRole(ctx context.Context, roleID int) ([]Permission, error)
		AddRolePermission(context.Context, RolePermission) (RolePermission, error)
		RemoveRolePermission(context.Context, RolePermission) error
	}

	// PermissionRepo defines all permission interactions
	PermissionRepo interface {
		ListAllPermissions(ctx context.Context) ([]Permission, error)
		ListPermissionsWithRolesCount(ctx context.Context) ([]PermissionWithRolesCount, error)
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
		LastLogin *time.Time `json:"lastLogin,omitempty"`
		Enabled   bool       `json:"enabled"`
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		CreatedBy *int64     `json:"createdBy,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy *int64     `json:"updatedBy,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy *int64     `json:"deletedBy,omitempty"`
		Roles     []Role     `json:"roles,omitempty"`
	}

	// UserFullDB represents a user and all columns
	UserFullDB struct {
		User
		LastLogin null.Time `db:"last_login"`
		Enabled   bool      `db:"enabled"`
		CreatedAt null.Time `db:"created_at"`
		CreatedBy null.Int  `db:"created_by"`
		UpdatedAt null.Time `db:"updated_at"`
		UpdatedBy null.Int  `db:"updated_by"`
		DeletedAt null.Time `db:"deleted_at"`
		DeletedBy null.Int  `db:"deleted_by"`
		Roles     []Role
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

	// RoleGetDTO represents relevant role fields for getting
	RoleGetDTO struct {
		RoleID int    `json:"roleID"`
		Name   string `json:"name"`
	}

	RoleAddEditDTO struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

	// RolePermission symbolises a link between a role.Role and permission.Permission
	RolePermission struct {
		RoleID       int `db:"role_id" json:"roleID"`
		PermissionID int `db:"permission_id" json:"permissionID"`
	}

	// RoleUser symbolises a link between a role.Role and User
	RoleUser struct {
		RoleID int `db:"role_id" json:"roleID"`
		UserID int `db:"user_id" json:"userID"`
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

	// Officership represents relevant officership fields
	Officership struct {
		OfficershipID    int     `json:"officershipID"`
		Name             string  `json:"name"`
		EmailAlias       string  `json:"emailAlias"`
		Description      string  `json:"description"`
		HistoryWikiURL   string  `json:"historyWikiURL"`
		RoleID           *int64  `json:"roleID,omitempty"`
		IsCurrent        bool    `json:"isCurrent"`
		IfUnfilled       *bool   `json:"ifUnfilled,omitempty"`
		CurrentOfficers  int     `json:"currentOfficers"`
		PreviousOfficers int     `json:"previousOfficers"`
		TeamID           *int64  `json:"teamID,omitempty"`
		TeamName         *string `json:"teamName,omitempty"`
		IsTeamLeader     *bool   `json:"isTeamLeader,omitempty"`
		IsTeamDeputy     *bool   `json:"isTeamDeputy,omitempty"`
	}

	// OfficershipGetDTO represents relevant officership fields for getting
	OfficershipGetDTO struct {
		OfficershipID int    `json:"officershipID"`
		Name          string `json:"name"`
	}

	// OfficershipAddEditDTO represents relevant officership fields for adding and editing
	OfficershipAddEditDTO struct {
		Name           string `json:"name"`
		EmailAlias     string `json:"emailAlias"`
		Description    string `json:"description"`
		HistoryWikiURL string `json:"historyWikiURL"`
		IsCurrent      bool   `json:"isCurrent"`
	}

	// OfficershipDB represents relevant officership fields
	OfficershipDB struct {
		OfficershipID    int         `db:"officer_id" json:"officershipID"`
		Name             string      `db:"name" json:"name"`
		EmailAlias       string      `db:"email_alias" json:"emailAlias"`
		Description      string      `db:"description" json:"description"`
		HistoryWikiURL   string      `db:"historywiki_url" json:"historyWikiURL"`
		RoleID           null.Int    `db:"role_id" json:"roleID,omitempty"`
		IsCurrent        bool        `db:"is_current" json:"isCurrent"`
		IfUnfilled       null.Bool   `db:"if_unfilled" json:"ifUnfilled,omitempty"`
		CurrentOfficers  int         `db:"current_officers" json:"currentOfficers,omitempty"`
		PreviousOfficers int         `db:"previous_officers" json:"previousOfficers,omitempty"`
		TeamID           null.Int    `db:"team_id" json:"teamID"`
		TeamName         null.String `db:"team_name" json:"teamName"`
		IsTeamLeader     null.Bool   `db:"is_team_leader" json:"isTeamLeader"`
		IsTeamDeputy     null.Bool   `db:"is_team_deputy" json:"isTeamDeputy"`
	}

	// OfficershipsStatus indicates the state desired for a database get of officers
	OfficershipsStatus int

	// OfficershipTeam represents relevant officership team fields
	//
	//nolint:revive
	OfficershipTeam struct {
		TeamID              int    `db:"team_id" json:"teamID"`
		Name                string `db:"name" json:"name"`
		EmailAlias          string `db:"email_alias" json:"emailAlias"`
		ShortDescription    string `db:"short_description" json:"shortDescription"`
		FullDescription     string `db:"full_description" json:"fullDescription"`
		CurrentOfficerships int    `db:"current_officerships" json:"currentOfficerships"`
		CurrentOfficers     int    `db:"current_officers" json:"currentOfficers"`
	}

	// OfficershipTeamAddEditDTO represents relevant officership team fields for adding and editing
	//
	//nolint:revive
	OfficershipTeamAddEditDTO struct {
		Name             string `db:"name" json:"name"`
		EmailAlias       string `db:"email_alias" json:"emailAlias"`
		ShortDescription string `db:"short_description" json:"shortDescription"`
		FullDescription  string `db:"full_description" json:"fullDescription"`
	}

	// OfficershipTeamGetDTO represents relevant officership team fields for getting
	//
	//nolint:revive
	OfficershipTeamGetDTO struct {
		TeamID int    `db:"team_id" json:"teamID"`
		Name   string `db:"name" json:"name"`
	}

	// OfficershipMember represents relevant officership member fields
	//
	//nolint:revive
	OfficershipMember struct {
		OfficershipMemberID int        `json:"officershipMemberID"`
		UserID              int        `json:"userID"`
		OfficerID           int        `json:"officerID"`
		StartDate           *time.Time `json:"startDate,omitempty"`
		EndDate             *time.Time `json:"endDate,omitempty"`
		OfficershipName     string     `json:"officershipName"`
		UserName            string     `json:"userName"`
		TeamID              *int       `json:"teamID,omitempty"`
		TeamName            *string    `json:"teamName,omitempty"`
	}

	// OfficershipMemberAddEditDTO represents relevant officership member fields
	//
	//nolint:revive
	OfficershipMemberAddEditDTO struct {
		UserID    int        `json:"userID"`
		OfficerID int        `json:"officerID"`
		StartDate *time.Time `json:"startDate,omitempty"`
		EndDate   *time.Time `json:"endDate,omitempty"`
	}

	// OfficershipMemberDB represents relevant officership member fields
	//
	//nolint:revive
	OfficershipMemberDB struct {
		OfficershipMemberID int         `db:"officership_member_id" json:"officershipMemberID"`
		UserID              int         `db:"user_id" json:"userID"`
		OfficerID           int         `db:"officer_id" json:"officerID"`
		StartDate           null.Time   `db:"start_date" json:"startDate"`
		EndDate             null.Time   `db:"end_date" json:"endDate"`
		OfficershipName     string      `db:"officership_name" json:"officershipName"`
		UserName            string      `db:"user_name" json:"userName"`
		TeamID              null.Int    `db:"team_id" json:"teamID"`
		TeamName            null.String `db:"team_name" json:"teamName"`
	}

	// OfficershipTeamMember represents relevant officership team member fields
	//
	//nolint:revive
	OfficershipTeamMember struct {
		TeamID           int    `db:"team_id" json:"officershipTeamMemberID"`
		OfficerID        int    `db:"officer_id" json:"officerID"`
		IsLeader         bool   `db:"is_leader" json:"isLeader"`
		IsDeputy         bool   `db:"is_deputy" json:"isDeputy"`
		IsCurrent        bool   `db:"is_current" json:"isCurrent"`
		OfficerName      string `db:"officer_name" json:"officerName"`
		CurrentOfficers  int    `db:"current_officers" json:"currentOfficers"`
		PreviousOfficers int    `db:"previous_officers" json:"previousOfficers"`
	}

	// OfficershipTeamMemberGetDeleteDTO represents relevant officership team member fields for getting and deleting
	//
	//nolint:revive
	OfficershipTeamMemberGetDeleteDTO struct {
		TeamID    int `db:"team_id" json:"officershipTeamMemberID"`
		OfficerID int `db:"officer_id" json:"officerID"`
	}

	// OfficershipTeamMemberAddDTO represents relevant officership team member fields for adding
	//
	//nolint:revive
	OfficershipTeamMemberAddDTO struct {
		TeamID    int  `db:"team_id" json:"officershipTeamMemberID"`
		OfficerID int  `db:"officer_id" json:"officerID"`
		IsLeader  bool `db:"is_leader" json:"isLeader"`
		IsDeputy  bool `db:"is_deputy" json:"isDeputy"`
	}

	CountOfficerships struct {
		TotalOfficerships   int `db:"total_officerships" json:"totalOfficerships"`
		CurrentOfficerships int `db:"current_officerships" json:"currentOfficerships"`
		TotalOfficers       int `db:"total_officers" json:"totalOfficers"`
		CurrentOfficers     int `db:"current_officers" json:"currentOfficers"`
	}
)

const (
	Any OfficershipsStatus = iota
	Retired
	Current
)

// NewStore creates a new store
func NewStore(db *sqlx.DB, cdn *s3.S3, cdnEndpoint string) Repo {
	return &Store{db: db, cdn: cdn, cdnEndpoint: cdnEndpoint}
}
