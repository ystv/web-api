package encoder

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/encoder"
	"github.com/ystv/web-api/utils"
)

type Repo struct {
	enc    *encoder.Encoder
	access *utils.Accesser
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

func NewEncoderController(enc *encoder.Encoder, access *utils.Accesser) *Repo {
	return &Repo{
		enc:    enc,
		access: access,
	}
}

// TODO: look into adding the parameter object without causing swagger to need to check external dependencies

// UploadRequest handles authenticating an upload request.
//
// Connects with tusd through web-hooks, so tusd POSTs here.
// Tusd's requests here do contain a lot of useful information.
// But for this endpoint, we are just checking for the JWT.
//
// @Summary New upload request
// @Description Authenticates tusd's webhook requests
// @ID new-encoder-upload-request
// @Tags encoder
// @Accept json
// @Success 200
// @Router /v1/internal/encoder/upload_request [post]
func (e *Repo) UploadRequest(c echo.Context) error {
	r := Request{}
	err := c.Bind(&r)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	_, err = e.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("GetToken failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

// TranscodeFinished handles marking a transcode item as finished
//
// @Summary Transcode Finished
// @Description Marks a transcode item as finished
// @ID new-encoder-transcode-finished
// @Tags encoder
// @Accept json
// @Param taskid path int true "Task ID"
// @Success 200
// @Router /v1/internal/encoder/transcode_finished/{taskid} [post]
func (e *Repo) TranscodeFinished(c echo.Context) error {
	err := e.enc.TranscodeFinished(c.Request().Context(), c.Param("taskid"))
	if err != nil {
		err = fmt.Errorf("transcode finished failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
