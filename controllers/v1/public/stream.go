package public

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils"
)

// ListChannels handles listing all channels
//
// @Summary Provides the visible channels
// @Description Lists the publicly visible channels
// @ID get-public-stream-channels
// @Tags public-playout-channels
// @Produce json
// @Success 200 {array} public.Channel
// @Router /v1/public/playout/channels [get]
func (s *Store) ListChannels(c echo.Context) error {
	chs, err := s.public.ListChannels(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public listchannels failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(chs))
}

// GetChannel handles listing teams and their members and info
//
// @Summary Provides a public or unlisted channel
// @ID get-public-stream-channel
// @Tags public-playout-channels
// @Param channelShortName path int true "Channel short name"
// @Produce json
// @Success 200 {object} public.Channel
// @Router /v1/public/playout/channel/{channelShortName} [get]
func (s *Store) GetChannel(c echo.Context) error {
	chs, err := s.public.GetChannel(c.Request().Context(), c.Param("channelShortName"))
	if err != nil {
		err = fmt.Errorf("public getchannel failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, chs)
}
