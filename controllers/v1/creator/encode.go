package creator

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/types/encode"
)

// EncodeProfileList handles listing encode formats
func (r *Repos) EncodeProfileList(c echo.Context) error {
	e, err := r.encode.ListFormat(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("EncodeProfileList failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, e)
}

// PresetList handles listing presets
func (r *Repos) PresetList(c echo.Context) error {
	p, err := r.encode.ListPreset(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("PresetList failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PresetNew handles creating a new preset
func (r *Repos) PresetNew(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	presetID, err := r.encode.NewPreset(c.Request().Context(), &p)
	if err != nil {
		err = fmt.Errorf("PresetNew failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, presetID)
}

// PresetUpdate handles updating a preset
func (r *Repos) PresetUpdate(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.encode.UpdatePreset(c.Request().Context(), &p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "No preset found")
		}
		err = fmt.Errorf("PresetUpdate failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
