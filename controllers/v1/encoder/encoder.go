package encoder

import (
	"log"
	"net/http"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
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
// Connects with tusd through web-hooks, so tusd POSTs here.
func VideoNew(c echo.Context) error {
	r := Request{}
	err := c.Bind(&r)
	log.Printf("%# v", pretty.Formatter(r))
	if err != nil {
		log.Print("VideoNew failed:")
		log.Printf("%# v", pretty.Formatter(err))
	}
	return c.JSON(http.StatusOK, nil)
}
