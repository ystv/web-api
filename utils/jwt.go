package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	// JWTClaims represents an identifiable JWT
	JWTClaims struct {
		UserID      int      `json:"id"`
		Permissions []string `json:"perms"`
		jwt.StandardClaims
	}
	// Permission represents the permissions that a user has
	Permission struct {
		PermissionID int    `json:"id"`
		Name         string `json:"name"`
	}
)

var (
	ErrNoToken      = errors.New("failed to find token")
	ErrInvalidToken = errors.New("invalid token")
)

// GetToken will return the JWT claims from a valid JWT token
func GetTokenEcho(c echo.Context) (*JWTClaims, error) {
	token := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer ")
	token = splitToken[1]

	if token == "" {
		return nil, ErrNoToken
	}
	return getClaims(token)
}

// GetToken will return the JWT claims from a valid JWT token
func GetTokenHTTP(r *http.Request) (*JWTClaims, error) {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer ")
	token = splitToken[1]

	if token == "" {
		return nil, ErrNoToken
	}
	return getClaims(token)
}

func getClaims(token string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("WAPI_SIGNING_KEY")), nil
	})
	if err != nil {
		err = fmt.Errorf("GetToken failed: %w", err)
		return nil, err
	}
	if claims.Valid() != nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
