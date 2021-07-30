package public

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ListChannels handles listing teams and their members and info
//
// @Summary Provides the visible channels
// @Description Lists the publically visible channels
// @ID get-public-stream-channels
// @Tags public-playout-channels
// @Produce json
// @Success 200 {array} public.Channel
// @Router /v1/public/playout/channels [get]
func (r *Repos) ListChannels(c echo.Context) error {
	chs, err := r.public.ListChannels(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("Public ListChannels failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, chs)
}
