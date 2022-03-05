package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type (
	Accesser struct {
		conf Config
	}
	Config struct {
		accessCookieName string
		signingKey       []byte
	}
	// AccessClaims represents an identifiable JWT
	AccessClaims struct {
		UserID      int      `json:"id"`
		Permissions []string `json:"perms"`
		jwt.StandardClaims
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

// GetToken will return the claims from a access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetToken(r *http.Request) (*AccessClaims, error) {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer ")
	token = splitToken[1]

	if token == "" {
		cookie, err := r.Cookie(a.conf.accessCookieName)
		token = cookie.Value
		if err != nil {
			return nil, fmt.Errorf("failed to get cookie", err)
		}
	}

	if token == "" {
		return nil, ErrNoToken
	}
	return a.getClaims(token)
}

func (a *Accesser) getClaims(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.conf.signingKey, nil
	})
	if err != nil {
		err = fmt.Errorf("failed to parse jwt: %w", err)
		return nil, err
	}
	if claims.Valid() != nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
