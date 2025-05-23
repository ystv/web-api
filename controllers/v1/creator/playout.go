package creator

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/creator/types/playout"
	"github.com/ystv/web-api/utils"
)

// ListChannels handles listing channels
// @Summary List all channels
// @Description Lists all channels, these are a rough implementation of what is to come (linear channels)
// @ID get-creator-playout-channels
// @Tags creator-playout-channels
// @Produce json
// @Success 200 {array} playout.Channel
// @Router /v1/internal/creator/playout/channels [get]
func (s *Store) ListChannels(c echo.Context) error {
	chs, err := s.channel.ListChannels(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListChannels failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(chs))
}

// NewChannel handles creating a new channel
// @Summary New channel
// @Description creates a new channel.
// @ID new-creator-playout-channel
// @Tags creator-playout-channels
// @Accept json
// @Param channel body playout.Channel true "Channel object"
// @Success 201 body int "Channel ID"
// @Router /v1/internal/creator/playout/channel [post]
func (s *Store) NewChannel(c echo.Context) error {
	var ch playout.Channel

	err := c.Bind(&ch)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.channel.NewChannel(c.Request().Context(), ch)
	if err != nil {
		err = fmt.Errorf("NewChannel failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusCreated)
}

// UpdateChannel handles updating a Channel
// @Summary Update a channel
// @Description updates a channel
// @ID update-creator-playout-channel
// @Tags creator-playout-channels
// @Accept json
// @Param channel body playout.Channel true "Channel object"
// @Success 200
// @Router /v1/internal/creator/playout/channel [put]
func (s *Store) UpdateChannel(c echo.Context) error {
	var ch playout.Channel

	err := c.Bind(&ch)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.channel.UpdateChannel(c.Request().Context(), ch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "No channel found")
		}
		err = fmt.Errorf("PresetChannel failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteChannel handles deleting channels
// @Summary Delete channel
// @Description deletes a channel by the short URL.
// @ID delete-creator-playout-channel
// @Tags creator-playout-channels
// @Param channelid path string true "Channel URL Name"
// @Success 200
// @Router /v1/internal/creator/playout/channel/{channelid} [delete]
func (s *Store) DeleteChannel(c echo.Context) error {
	err := s.channel.DeleteChannel(c.Request().Context(), c.Param("channelid"))
	if err != nil {
		err = fmt.Errorf("DeleteChannel failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
