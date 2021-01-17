package utils

import "github.com/dgrijalva/jwt-go"

// JWTClaims represents an identifiable JWT
type JWTClaims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}
