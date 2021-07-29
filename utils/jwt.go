package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	// JWTClaims represents an identifiable JWT
	JWTClaims struct {
		UserID      int          `json:"id"`
		Permissions []Permission `json:"perms"`
		jwt.StandardClaims
	}
	// Permission represents the permissions that a user has
	Permission struct {
		PermissionID int    `json:"id"`
		Name         string `json:"name"`
	}
)

var (
	ErrNoCookie      = errors.New("failed to find token cookie")
	ErrInvalidCookie = errors.New("invalid cookie")
)

func GetToken(cookie *http.Cookie) (*JWTClaims, error) {
	tokenString := cookie.Value
	claims := &JWTClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("WAPI_SIGNING_KEY")), nil
	})
	if err != nil {
		err = fmt.Errorf("GetToken failed: %w", err)
		return nil, err
	}
	if claims.Valid() != nil {
		return nil, ErrInvalidCookie
	}
	return claims, nil
}

// GetToken will return the JWT claims from a valid JWT token
func GetTokenEcho(c echo.Context) (*JWTClaims, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return nil, ErrNoCookie
	}
	return GetToken(cookie)
}

// GetToken will return the JWT claims from a valid JWT token
func GetTokenHTTP(r *http.Request) (*JWTClaims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, ErrNoCookie
	}
	return GetToken(cookie)
}
