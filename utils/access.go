package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils/permissions/users"
)

type (
	Repo interface {
		GetToken(r *http.Request) (*AccessClaims, int, error)
		AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		AddUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		ListUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		GroupAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		PermissionsAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		OfficershipAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		SuperUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		ModifyUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		ManageStreamAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	}

	Accesser struct {
		conf Config
	}

	Config struct {
		AccessCookieName string
		SigningKey       []byte
	}

	// AccessClaims represents an identifiable JWT
	AccessClaims struct {
		UserID      int      `json:"id"`
		Permissions []string `json:"perms"`
		jwt.RegisteredClaims
	}

	// Permission represents the permissions that a user has
	Permission struct {
		Name string `json:"name"`
	}
)

var (
	ErrNoToken      = errors.New("token not found")
	ErrInvalidToken = errors.New("invalid token")
)

// NewAccesser allows the validation of web-auth JWT tokens both as
// headers and as cookies
func NewAccesser(conf Config) Repo {
	return &Accesser{
		conf: conf,
	}
}

// GetToken will return the claims from an access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetToken(r *http.Request) (*AccessClaims, int, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.AccessCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return nil, http.StatusUnauthorized, ErrNoToken
			}
			return nil, http.StatusBadRequest, fmt.Errorf("failed to get cookie: %w", err)
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return nil, http.StatusBadRequest, ErrInvalidToken
		}
		token = splitToken[1]
	}

	if token == "" {
		return nil, http.StatusUnauthorized, ErrNoToken
	}

	claims := &AccessClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		log.Printf("error with signing: %+v", err)
		return nil, http.StatusUnauthorized, ErrInvalidToken
	}

	return claims, http.StatusOK, nil
}

// AuthMiddleware checks an HTTP request for a valid token either in the header or cookie
func (a *Accesser) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		return next(c)
	}
}

// AddUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie,
// and if the user can add a user
func (a *Accesser) AddUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersMembersAdmin || p == users.ManageMembersMembersAdd {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// ListUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can list users
func (a *Accesser) ListUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersMembersAdmin || p == users.ManageMembersMembersList {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// GroupAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can modify groups
func (a *Accesser) GroupAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersGroup {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// PermissionsAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can modify permissions
func (a *Accesser) PermissionsAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersPermissions {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// SuperUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie
// and if the user is SuperUser
func (a *Accesser) SuperUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// OfficershipAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can modify permissions
func (a *Accesser) OfficershipAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersOfficers {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// ModifyUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can list users
func (a *Accesser) ModifyUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersMembersAdmin {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}

// ManageStreamAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can manage streams
func (a *Accesser) ManageStreamAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, status, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     status,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.Cobra {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
}
