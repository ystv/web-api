package customsettings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/customsettings"
	"github.com/ystv/web-api/utils"
)

type (
	Repos interface {
		SettingsRepo
	}

	SettingsRepo interface {
		ListCustomSettings(c echo.Context) error
		GetCustomSetting(c echo.Context) error
		AddCustomSetting(c echo.Context) error
		EditCustomSetting(c echo.Context) error
		DeleteCustomSetting(c echo.Context) error
	}

	// Store stores our dependencies
	Store struct {
		access         utils.Repo
		customSettings customsettings.Repo
	}
)

// NewRepos creates our data store
func NewRepos(db *sqlx.DB, access utils.Repo) Repos {
	return &Store{
		access:         access,
		customSettings: customsettings.NewStore(db),
	}
}

// ListCustomSettings handles listing settings
//
// @Summary Provides all settings
// @Description Contains settings and value in string base64 format
// @ID get-custom-settings
// @Tags custom-settings
// @Produce json
// @Success 200 {array} customsettings.CustomSetting
// @Router /v1/internal/custom-settings [get]
func (s *Store) ListCustomSettings(c echo.Context) error {
	settings, err := s.customSettings.ListCustomSettings(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public ListCustomSettings failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, settings)
}

// GetCustomSetting handles getting a single setting
//
// @Summary Provides a single setting
// @Description Contains a single setting value in string base64 format
// @ID get-custom-setting
// @Tags custom-settings
// @Param settingid path string true "Setting id"
// @Produce json
// @Success 200 {object} customsettings.CustomSetting
// @Router /v1/internal/custom-setting/{settingid} [get]
func (s *Store) GetCustomSetting(c echo.Context) error {
	setting, err := s.customSettings.GetCustomSetting(c.Request().Context(), c.Param("settingid"))
	if err != nil {
		err = fmt.Errorf("public GetTeamByEmail failed: %w", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, setting)
}

// AddCustomSetting handles creating a custom setting
//
// @Summary Create a custom setting
// @ID add-custom-setting
// @Tags custom-settings
// @Produce json
// @Param customSetting body customsettings.CustomSetting true "Custom Setting object"
// @Success 201 {object} customsettings.CustomSetting
// @Router /v1/internal/custom-setting [post]
func (s *Store) AddCustomSetting(c echo.Context) error {
	var customSettingAdd customsettings.CustomSetting
	err := c.Bind(&customSettingAdd)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("request body could not be decoded: %w", err))
	}

	if customSettingAdd.SettingID == "" || customSettingAdd.Value == "" {
		fmt.Println(customSettingAdd.SettingID, customSettingAdd.Value)
		return echo.NewHTTPError(http.StatusBadRequest, "Setting ID and value must be filled for add custom setting")
	}

	cs1, err := s.customSettings.GetCustomSetting(c.Request().Context(), customSettingAdd.SettingID)
	if err == nil && len(cs1.SettingID) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "custom setting with setting ID \""+customSettingAdd.SettingID+"\" already exists")
	}

	if val, ok := customSettingAdd.Value.(string); ok {
		if !json.Valid([]byte(val)) {
			return echo.NewHTTPError(http.StatusBadRequest, "custom setting value is not json and invalid")
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "custom setting value is invalid")
	}

	customSetting, err := s.customSettings.AddCustomSetting(c.Request().Context(), customSettingAdd)
	if err != nil {
		return fmt.Errorf("failed to add custom setting for add custom setting: %w", err)
	}

	return c.JSON(http.StatusCreated, customSetting)
}

// EditCustomSetting handles editing a custom setting
//
// @Summary Edits a custom setting
// @ID edit-custom-setting
// @Tags custom-settings
// @Produce json
// @Param settingid path int true "setting id"
// @Param customSetting body customsettings.CustomSettingEditDTO true "Custom Setting object"
// @Success 200 {object} customsettings.CustomSetting
// @Router /v1/internal/custom-setting/{officerid} [put]
func (s *Store) EditCustomSetting(c echo.Context) error {
	settingID := c.Param("settingid")

	_, err := s.customSettings.GetCustomSetting(c.Request().Context(), settingID)
	if err != nil {
		return fmt.Errorf("failed to get custom setting for edit custom setting: %w", err)
	}

	var customSettingEdit customsettings.CustomSettingEditDTO
	err = c.Bind(&customSettingEdit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("request body could not be decoded: %w", err))
	}

	if customSettingEdit.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "value must be filled for edit custom setting")
	}

	if !json.Valid([]byte(customSettingEdit.Value)) {
		return echo.NewHTTPError(http.StatusBadRequest, "value is not json and invalid")
	}

	customSetting, err := s.customSettings.EditCustomSetting(c.Request().Context(), settingID, customSettingEdit)
	if err != nil {
		return fmt.Errorf("failed to edit custom setting for edit custom setting: %w", err)
	}

	return c.JSON(http.StatusOK, customSetting)
}

// DeleteCustomSetting handles deleting a custom setting
//
// @Summary Deletes custom setting
// @ID delete-custom-setting
// @Tags custom-settings
// @Param settingid path int true "setting id"
// @Produce json
// @Success 204
// @Router /v1/internal/custom-setting/{settingid} [delete]
func (s *Store) DeleteCustomSetting(c echo.Context) error {
	customSetting, err := s.customSettings.GetCustomSetting(c.Request().Context(), c.Param("settingid"))
	if err != nil {
		return fmt.Errorf("failed to get custom setting for custom setting delete: %w", err)
	}

	err = s.customSettings.DeleteCustomSetting(c.Request().Context(), customSetting.SettingID)
	if err != nil {
		return fmt.Errorf("failed to delete custom setting for custom setting delete: %w", err)
	}

	return c.NoContent(http.StatusNoContent)
}
