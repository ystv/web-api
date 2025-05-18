package public

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCustomSettingPublic handles getting a single setting value that is public
//
// @Summary Provides a single public setting
// @Description Contains a single setting value in string format
// @ID get-public-custom-setting
// @Tags public-custom-settings
// @Param settingid path string true "Setting id"
// @Produce json
// @Success 200 {object} public.CustomSetting
// @Router /v1/public/custom-setting/{settingid} [get]
func (s *Store) GetCustomSettingPublic(c echo.Context) error {
	setting, err := s.public.GetCustomSettingPublic(c.Request().Context(), c.Param("settingid"))
	if err != nil {
		err = fmt.Errorf("public GetCustomSettingPublic failed: %w", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, setting)
}
