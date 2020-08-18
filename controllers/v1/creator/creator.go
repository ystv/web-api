package creator

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator"
)

// FileUpload Handles uploading a file
func FileUpload(c echo.Context) error {
	creator.CreateBucket("pending", "ystv-wales-1")
	url, err := creator.GenerateUploadURL("pending", c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, url)
}

// Stats handles sending general stats about the video library
func Stats(c echo.Context) error {
	s, err := creator.Stats(context.Background())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, s)
}
