package creator

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/types/encode"
)

// EncodeProfileList handles listing encode formats
func (r *Repos) EncodeProfileList(c echo.Context) error {
	e, err := r.encode.ListFormat(c.Request().Context())
	if err != nil {
		log.Printf("EncodeProfileList failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, e)
}

// PresetList handles listing presets
func (r *Repos) PresetList(c echo.Context) error {
	p, err := r.encode.ListPreset(c.Request().Context())
	if err != nil {
		log.Printf("PresetList failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PresetNew handles creating a new preset
func (r *Repos) PresetNew(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	presetID, err := r.encode.NewPreset(c.Request().Context(), &p)
	if err != nil {
		log.Printf("PresetNew failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, presetID)
}

// PresetUpdate handles updating a preset
func (r *Repos) PresetUpdate(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = r.encode.UpdatePreset(c.Request().Context(), &p)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, err)
		}
		log.Printf("PresetUpdate failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
