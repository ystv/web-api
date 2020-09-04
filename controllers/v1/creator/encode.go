package creator

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/types/encode"
)

// ListEncodeProfile handles listing encode formats
// @Summary List all encode formats
// @Description Lists all encode formats, these are instructions for the encoder to create the video
// @ID get-creator-encodes-formats
// @Tags creator, encodes
// @Produce json
// @Success 200 {array} encode.Format
// @Router /v1/internal/creator/encodes/profiles [get]
func (r *Repos) ListEncodeProfile(c echo.Context) error {
	e, err := r.encode.ListFormat(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("EncodeProfileList failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, e)
}

// ListPreset handles listing presets
// @Summary List all encode presets
// @Description Lists all encode presets, these are groups of instructions (formats) for the encoder to create the video
// @ID get-creator-encodes-presets
// @Tags creator, encodes
// @Produce json
// @Success 200 {array} encode.Preset
// @Router /v1/internal/creator/encodes/presets [get]
func (r *Repos) ListPreset(c echo.Context) error {
	p, err := r.encode.ListPreset(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("PresetList failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// NewPreset handles creating a new preset
// @Summary New preset
// @Description creates a new preset.
// @ID new-creator-encodes-preset
// @Tags creator, encodes
// @Accept json
// @Param event body encode.Preset true "Preset object"
// @Success 201 body int "Preset ID"
// @Router /v1/internal/creator/encodes/presets [post]
func (r *Repos) NewPreset(c echo.Context) error {
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

// UpdatePreset handles updating a preset
// @Summary Update a preset
// @Description updates an preset
// @ID update-creator-encodes-preset
// @Tags creator, encodes
// @Accept json
// @Param quote body encode.Preset true "Preset object"
// @Success 200
// @Router /v1/internal/creator/encodes/presets [put]
func (r *Repos) UpdatePreset(c echo.Context) error {
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
