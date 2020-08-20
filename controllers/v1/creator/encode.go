package creator

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/encode"
)

// EncodeProfileList handles listing encode formats
func EncodeProfileList(c echo.Context) error {
	e, err := encode.FormatList()
	if err != nil {
		log.Printf("EncodeProfileList failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, e)
}

// PresetList handles listing presets
func PresetList(c echo.Context) error {
	p, err := encode.PresetList()
	if err != nil {
		log.Printf("PresetList failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PresetNew handles creating a new preset
func PresetNew(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = encode.PresetNew(&p)
	if err != nil {
		log.Printf("PresetNew failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// PresetUpdate handles updating a preset
func PresetUpdate(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = encode.PresetUpdate(&p)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, err)
		}
		log.Printf("PresetUpdate failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
