package encoder

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
)

type (
	// These structs are for binding to tusd's request

	// Request represents the upload and a normal HTTP request
	Request struct {
		Upload      Upload
		HTTPRequest http.Request
	}
	// Upload represents an object and it's status
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
	// MetaData represents metadata of a file.
	// There is more, but we just need filename
	MetaData struct {
		Filename string `json:"filename"`
	}
	// Storage represents the storage medium of the object
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
		err = fmt.Errorf("VideoNew failed: failed to find api token: %w", err)
		return echo.NewHTTPError(http.StatusForbidden, err)
	}
	tokenString := cookie.Value
	claims := &people.JWTClaims{}

	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("signing_key")), nil
	})
	if err != nil {
		err = fmt.Errorf("VideoNew failed: failed to parse jwt %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if claims.Valid() != nil {
		err = fmt.Errorf("JWT expired: %w", err)
		return echo.NewHTTPError(http.StatusForbidden, err)
	}
	return c.NoContent(http.StatusOK)
}
