package encoder

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
)

type (
	Request struct {
		Upload      Upload
		HTTPRequest http.Request
	}
	Upload struct {
		ID        string
		Size      int
		Offset    int
		IsFinal   bool
		IsPartial bool
		// PartialUploads null
		MetaData []MetaData
		Storage  Storage
	}
	MetaData struct {
		Filename string `json:"filename"`
	}
	Storage struct {
		Type   string
		Bucket string
		Key    string
	}
)

// VideoNew handles authenticating a video upload request.
//
// Connects with tusd through web-hooks, so tusd POSTs here.
// tusd's requests here does contain a lot of useful information.
// but for this endpoint, we are just checking for the JWT.
func VideoNew(c echo.Context) error {
	r := Request{}
	c.Bind(&r)
	if r.HTTPRequest.Method != "POST" {
		return c.NoContent(http.StatusOK)
	}
	cookie, err := r.HTTPRequest.Cookie("token")
	if err != nil {
		return echo.ErrBadRequest
	}
	tokenString := cookie.Value
	claims := &people.JWTClaims{}

	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("signing_key")), nil
	})
	if err != nil {
		log.Printf("UserByToken failed: %+v", err)
		return echo.ErrInternalServerError
	}
	if claims.Valid() != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, nil)
}