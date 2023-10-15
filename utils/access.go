package utils

import (
	"errors"
	"fmt"
	"github.com/ystv/web-api/utils/permissions/users"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	Accesser struct {
		conf Config
	}

	Config struct {
		AccessCookieName string
		SecurityBaseURL  string
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
func NewAccesser(conf Config) *Accesser {
	return &Accesser{
		conf: conf,
	}
}

// GetToken will return the claims from an access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetToken(r *http.Request) (*AccessClaims, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.AccessCookieName)
		if err != nil {
			if errors.As(http.ErrNoCookie, &err) {
				return nil, ErrNoToken
			}
			return nil, fmt.Errorf("failed to get cookie: %w", err)
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return nil, ErrInvalidToken
		}
		token = splitToken[1]
	}

	if token == "" {
		return nil, ErrNoToken
	}
	return a.getClaims(token)
}

func (a *Accesser) getClaims(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// AuthMiddleware checks an HTTP request for a valid token either in the header or cookie
func (a *Accesser) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  err.Error(),
				Internal: err,
			}
		}
		return next(c)
	}
}

// AddUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can add a user
func (a *Accesser) AddUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersMembersAdmin || p == users.ManageMembersMembersAdd {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
}

// ListUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can list users
func (a *Accesser) ListUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersMembersAdmin || p == users.ManageMembersMembersList {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
}

// ModifyUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can list users
func (a *Accesser) ModifyUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  err.Error(),
				Internal: err,
			}
		}
		for _, p := range claims.Permissions {
			if p == users.SuperUser || p == users.ManageMembersAdmin || p == users.ManageMembersMembersAdmin {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
}
