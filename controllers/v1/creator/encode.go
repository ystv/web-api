package creator

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/creator/types/encode"
	"github.com/ystv/web-api/utils"
)

// ListEncodeFormat handles listing encode formats
// @Summary List all encode formats
// @Description Lists all encode formats, these are instructions for the encoder to create the video
// @ID get-creator-encode-format
// @Tags creator-encodes
// @Produce json
// @Success 200 {array} encode.Format
// @Router /v1/internal/creator/encode/format [get]
func (r *Repos) ListEncodeFormat(c echo.Context) error {
	e, err := r.encode.ListFormat(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListFormat failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(e))
}

// NewEncodeFormat handles creating a new encode format
// @Summary New encode format
// @Description creates a new encode format.
// @ID new-creator-encode-format
// @Tags creator-encodes
// @Accept json
// @Param format body encode.Format true "Encode format object"
// @Success 201 body int "Format ID"
// @Router /v1/internal/creator/encode/format [post]
func (r *Repos) NewEncodeFormat(c echo.Context) error {
	format := encode.Format{}
	err := c.Bind(&format)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	formatID, err := r.encode.NewFormat(c.Request().Context(), format)
	if err != nil {
		err = fmt.Errorf("NewFormat failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, formatID)
}

// UpdateEncodeFormat handles updating a format
// @Summary Update a format
// @Description updates a format
// @ID update-creator-encode-format
// @Tags creator-encodes
// @Accept json
// @Param format body encode.Format true "Format object"
// @Success 200
// @Router /v1/internal/creator/encode/format [put]
func (r *Repos) UpdateEncodeFormat(c echo.Context) error {
	format := encode.Format{}
	err := c.Bind(&format)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.encode.UpdateFormat(c.Request().Context(), format)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "no preset found")
		}
		err = fmt.Errorf("PresetUpdate failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// DeleteEncodeFormat handles deleting quotes
// @Summary Delete an encode format
// @Description Delete a video encode format
// @ID delete-creator-encode-format
// @Tags creator-encodes
// @Param formatid path int true "Format ID"
// @Success 200
// @Router /v1/internal/creator/encode/format/{formatid} [delete]
func (r *Repos) DeleteEncodeFormat(c echo.Context) error {
	formatID, err := strconv.Atoi(c.Param("formatid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	err = r.encode.DeleteFormat(c.Request().Context(), formatID)
	if err != nil {
		err = fmt.Errorf("DeleteFormat failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// ListEncodePreset handles listing presets
// @Summary List all encode presets
// @Description Lists all encode presets, these are groups of instructions (formats) for the encoder to create the video
// @ID get-creator-encode-preset
// @Tags creator-encodes
// @Produce json
// @Success 200 {array} encode.Preset
// @Router /v1/internal/creator/encode/preset [get]
func (r *Repos) ListEncodePreset(c echo.Context) error {
	p, err := r.encode.ListPreset(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListPreset failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// NewEncodePreset handles creating a new preset
// @Summary New preset
// @Description creates a new preset.
// @ID new-creator-encode-preset
// @Tags creator-encodes
// @Accept json
// @Param event body encode.Preset true "Preset object"
// @Success 201 body int "Preset ID"
// @Router /v1/internal/creator/encode/preset [post]
func (r *Repos) NewEncodePreset(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	presetID, err := r.encode.NewPreset(c.Request().Context(), p)
	if err != nil {
		err = fmt.Errorf("PresetNew failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, presetID)
}

// UpdateEncodePreset handles updating a preset
// @Summary Update a preset
// @Description updates a preset
// @ID update-creator-encode-preset
// @Tags creator-encodes
// @Accept json
// @Param quote body encode.Preset true "Preset object"
// @Success 200
// @Router /v1/internal/creator/encode/preset [put]
func (r *Repos) UpdateEncodePreset(c echo.Context) error {
	p := encode.Preset{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.encode.UpdatePreset(c.Request().Context(), p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "no preset found")
		}
		err = fmt.Errorf("PresetUpdate failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// DeleteEncodePreset handles deleting presets
// @Summary Delete a encode preset
// @Description Delete a video encode preset
// @ID delete-creator-encode-preset
// @Tags creator-encodes
// @Param presetid path int true "Preset ID"
// @Success 200
// @Router /v1/internal/creator/encode/preset/{presetid} [delete]
func (r *Repos) DeleteEncodePreset(c echo.Context) error {
	presetID, err := strconv.Atoi(c.Param("presetid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	err = r.encode.DeletePreset(c.Request().Context(), presetID)
	if err != nil {
		err = fmt.Errorf("DeletePreset failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
