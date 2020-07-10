package creator

import (
	"net/http"

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

// CreationFind Handles finding a creation by ID
func CreationFind(c echo.Context) error {
	creation, _ := creator.VideoItemFind()
	return c.JSON(http.StatusOK, creation)
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
