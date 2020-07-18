package creator

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator"
)

// CreationMetaCreate Handes uploading meta data for a creation
func CreationMetaCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Meta created")
}

// CreationFileUpload Handles uploading a file
func CreationFileUpload(c echo.Context) error {
	creator.CreateBucket("pending", "ystv-wales-1")
	url, err := creator.GenerateUploadURL("pending", c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, url)
}

// VideoFind finds a video by ID
func VideoFind(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	v, err := creator.VideoItemFind(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
}

// CreationCreate Handles creation of a creation lol
func CreationCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Creation created")
}

// CreationList Handles listing all creations
func CreationList(c echo.Context) error {
	creations, err := creator.ListPendingUploads()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, creations)
}
