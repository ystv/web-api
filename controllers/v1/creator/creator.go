package creator

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services"
)

// CreationMetaCreate Handes uploading meta data for a creation
func CreationMetaCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Meta created")
}

// CreationFileUpload Handles uploading a file
func CreationFileUpload(c echo.Context) error {
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	return err
	// }
	// src, err := file.Open()
	// if err != nil {
	// 	return err
	// }
	// defer src.Close()

	// dst, err := os.Create(file.Filename)
	// if err != nil {
	// 	return err
	// }
	// defer dst.Close()

	// if _, err = io.Copy(dst, src); err != nil {
	// 	return err
	// }
	// return c.String(http.StatusOK, "Lovely")
	services.CreateBucket("pending", "ystv-wales-1")
	url, err := services.GenerateUploadURL("pending", c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, url)
}

// CreationFind Handles finding a creation by ID
func CreationFind(c echo.Context) error {
	return c.String(http.StatusOK, "Found creation")
}

// CreationCreate Handles creation of a creation lol
func CreationCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Creation created")
}

// CreationList Handles listing all creations
func CreationList(c echo.Context) error {
	creations, err := services.ListPendingUploads()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, creations)
}
