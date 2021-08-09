package encoder

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/encoder"
	"github.com/ystv/web-api/utils"
)

type EncoderController struct {
	enc *encoder.Encoder
}

type (
	// These structs are for binding to tusd's request

	// Request represents the upload and a normal HTTP request
	Request struct {
		Upload      Upload
		HTTPRequest *http.Request
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
func (e *EncoderController) VideoNew(c echo.Context) error {
	r := Request{}
	c.Bind(&r)
	if r.HTTPRequest.Method != "POST" {
		return c.NoContent(http.StatusOK)
	}

	_, err := utils.GetTokenHTTP(r.HTTPRequest)
	if err != nil {
		err = fmt.Errorf("GetToken failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

func (e *EncoderController) TranscodeFinished(c echo.Context) error {
	err := e.enc.TranscodeFinished(c.Request().Context(), c.Param("taskid"))
	if err != nil {
		err = fmt.Errorf("transcode finished failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

func NewEncoderController(enc *encoder.Encoder) *EncoderController {
	return &EncoderController{enc: enc}
}
