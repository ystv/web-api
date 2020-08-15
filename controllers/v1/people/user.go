package people

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/people"
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

// UserByID finds a user by ID
func UserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	p, err := people.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// UserByIDFull finds a user by ID returing all info
func UserByIDFull(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	p, err := people.GetFull(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// UserByToken finds a user by their JWT token
func UserByToken(c echo.Context) error {
	claims, err := GetToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	p, err := people.Get(claims.UserID)
	if err != nil {
		log.Printf("UserByToken failed getting: %+v", err)
	}
	return c.JSON(http.StatusOK, p)
}

// UserByTokenFull finds a user by their JWT token returning all info
func UserByTokenFull(c echo.Context) error {
	claims, err := GetToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	p, err := people.GetFull(claims.UserID)
	if err != nil {
		log.Printf("UserByToken failed getting: %+v", err)
	}
	return c.JSON(http.StatusOK, p)
}

// GetToken will return the JWT claims from a valid JWT token
func GetToken(c echo.Context) (*JWTClaims, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return nil, echo.ErrBadRequest
	}
	tokenString := cookie.Value
	claims := &JWTClaims{}

	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("signing_key")), nil
	})
	if err != nil {
		log.Printf("UserByToken failed: %+v", err)
		return nil, echo.ErrInternalServerError
	}
	return claims, nil
}
