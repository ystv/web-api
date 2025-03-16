package people

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

type (
	Repos interface {
		PermissionRepo
		RoleRepo
		UserRepo
		OfficershipRepo
	}
	OfficershipRepo interface {
		OfficershipCount(c echo.Context) error
		ListOfficerships(c echo.Context) error
		GetOfficership(c echo.Context) error
		AddOfficership(c echo.Context) error
		EditOfficership(c echo.Context) error
		DeleteOfficership(c echo.Context) error
		ListOfficershipTeams(c echo.Context) error
		GetOfficershipTeam(c echo.Context) error
		AddOfficershipTeam(c echo.Context) error
		EditOfficershipTeam(c echo.Context) error
		DeleteOfficershipTeam(c echo.Context) error
		ListOfficers(c echo.Context) error
		GetOfficer(c echo.Context) error
		AddOfficer(c echo.Context) error
		EditOfficer(c echo.Context) error
		DeleteOfficer(c echo.Context) error
		OfficershipTeamAddOfficership(c echo.Context) error
		OfficershipTeamRemoveOfficership(c echo.Context) error
	}

	PermissionRepo interface {
		ListPermissions(c echo.Context) error
		ListPermissionsWithRolesCount(c echo.Context) error
		ListPermissionMembersByID(c echo.Context) error
		GetPermissionByID(c echo.Context) error
		GetPermissionByIDWithRolesCount(ctx echo.Context) error
		AddPermission(c echo.Context) error
		EditPermission(c echo.Context) error
		DeletePermission(c echo.Context) error
	}

	RoleRepo interface {
		ListAllRolesWithPermissions(c echo.Context) error
		ListAllRolesWithCount(c echo.Context) error
		GetRoleFull(c echo.Context) error
		ListRoleMembersByID(c echo.Context) error
		ListRolePermissionsByID(c echo.Context) error
		AddRole(c echo.Context) error
		EditRole(c echo.Context) error
		DeleteRole(c echo.Context) error
		RoleAddPermissionFunc(c echo.Context) error
		RoleRemovePermissionFunc(c echo.Context) error
		RoleAddUserFunc(c echo.Context) error
		RoleRemoveUserFunc(c echo.Context) error
	}

	UserRepo interface {
		UserStats(c echo.Context) error
		UserByID(c echo.Context) error
		UserByIDFull(c echo.Context) error
		UserByEmail(c echo.Context) error
		UserByEmailFull(c echo.Context) error
		UserByToken(c echo.Context) error
		UserByTokenFull(c echo.Context) error
		AddUser(c echo.Context) error
		ListAllPeople(c echo.Context) error
		ListPeoplePagination(c echo.Context) error
	}

	// Store stores our dependencies
	Store struct {
		people people.Repo
		access utils.Repo
	}
)

// NewRepos creates our data store
func NewRepos(db *sqlx.DB, cdn *s3.S3, access utils.Repo, cdnEndpoint string) Repos {
	return &Store{
		people: people.NewStore(db, cdn, cdnEndpoint),
		access: access,
	}
}

func (s *Store) UserFullDBToUserFull(userFullDB people.UserFullDB) people.UserFull {
	var lastLogin, createdAt, updatedAt, deletedAt *time.Time
	var createdBy, updatedBy, deletedBy *int64
	if userFullDB.LastLogin.Valid {
		lastLogin = &userFullDB.LastLogin.Time
	}
	if userFullDB.CreatedAt.Valid {
		createdAt = &userFullDB.CreatedAt.Time
	}
	if userFullDB.UpdatedAt.Valid {
		updatedAt = &userFullDB.UpdatedAt.Time
	}
	if userFullDB.DeletedAt.Valid {
		deletedAt = &userFullDB.DeletedAt.Time
	}
	if userFullDB.CreatedBy.Valid {
		createdBy = &userFullDB.CreatedBy.Int64
	}
	if userFullDB.UpdatedBy.Valid {
		updatedBy = &userFullDB.UpdatedBy.Int64
	}
	if userFullDB.DeletedBy.Valid {
		deletedBy = &userFullDB.DeletedBy.Int64
	}

	return people.UserFull{
		User:      userFullDB.User,
		LastLogin: lastLogin,
		Enabled:   userFullDB.Enabled,
		CreatedAt: createdAt,
		CreatedBy: createdBy,
		UpdatedAt: updatedAt,
		UpdatedBy: updatedBy,
		DeletedAt: deletedAt,
		DeletedBy: deletedBy,
		Roles:     userFullDB.Roles,
	}
}
