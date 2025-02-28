package people

import (
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
	}

	PermissionRepo interface {
		ListAllPermissions(c echo.Context) error
		ListPermissionMembersByID(c echo.Context) error
	}

	RoleRepo interface {
		ListAllRolesWithPermissions(c echo.Context) error
		ListAllRolesWithCount(c echo.Context) error
		ListRoleMembersByID(c echo.Context) error
		ListRolePermissionsByID(c echo.Context) error
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
