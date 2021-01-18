package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetPlaylist handles returning a playlist with a list of videos and metadata
//
// @Summary Provides a playlist
// @Description Returns a playlist object, includes videos (not videofiles) and metadata.
// @ID get-public-playlist
// @Tags public, playlist
// @Param playlistid path int true "Playlist ID"
// @Produce json
// @Success 200 {object} public.Playlist
// @Router /v1/public/playlist/{seriesid} [get]
func (r *Repos) GetPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlistid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad playlist ID")
	}
	s, err := r.public.GetPlaylist(c.Request().Context(), playlistID)
	if err != nil {
		err = fmt.Errorf("Public GetPlaylist failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, s)
}
