package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils/permissions/users"
)

type (
	Repo interface {
		GetToken(r *http.Request) (*AccessClaims, error)
		AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		AddUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		ListUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		ModifyUserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
		ManageStreamAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	}

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
func NewAccesser(conf Config) Repo {
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
			if errors.Is(err, http.ErrNoCookie) {
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
	return a.getClaims(r.Context(), token)
}

func (a *Accesser) getClaims(ctx context.Context, token string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		log.Printf("error with signing: %+v", err)
		return nil, ErrInvalidToken
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", a.conf.SecurityBaseURL+"/api/test", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		var b []byte

		b, err = io.ReadAll(res.Body)
		if err != nil {
			log.Printf("converting fail: %+v", err)
		}

		return nil, fmt.Errorf("invalid token: invalid status: %s: %s", res.Status, string(b))
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

// AddUserAuthMiddleware checks an HTTP request for a valid token either in the header or cookie,
// and if the user can add a user
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

// ManageStreamAuthMiddleware checks an HTTP request for a valid token either in the header or cookie and if the user can manage streams
func (a *Accesser) ManageStreamAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
			if p == users.SuperUser || p == users.Cobra {
				return next(c)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
}
