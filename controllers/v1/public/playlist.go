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
// @Description Returns a playlist object, includes videos (not video files) and metadata.
// @ID get-public-playlist
// @Tags public-playlist
// @Param playlistid path int true "Playlist ID"
// @Produce json
// @Success 200 {object} public.Playlist
// @Router /v1/public/playlist/{seriesid} [get]
func (s *Store) GetPlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlistid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad playlist ID")
	}

	p, err := s.public.GetPlaylist(c.Request().Context(), playlistID)
	if err != nil {
		err = fmt.Errorf("public GetPlaylist failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// GetPlaylistPopularByAllTime returns a fake playlist with a list of popular videos of all time
//
// @Summary Provides a playlist of popular videos of all time
// @Description Provides a fake playlist, containing a list of popular videos
// @ID get-public-playlist-popular-by-all-time
// @Tags public-playlist
// @Produce json
// @Success 200 {object} public.Playlist
// @Router /v1/public/playlist/popular/all [get]
func (s *Store) GetPlaylistPopularByAllTime(c echo.Context) error {
	p, err := s.public.GetPlaylistPopularByAllTime(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public GetPlaylistByAllTime failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// GetPlaylistPopularByPastYear returns a fake playlist with a list of popular videos of past year
//
// @Summary Provides a playlist of popular videos of past year
// @Description Provides a fake playlist, containing a list of popular videos
// @ID get-public-playlist-popular-by-past-year
// @Tags public-playlist
// @Produce json
// @Success 200 {object} public.Playlist
// @Router /v1/public/playlist/popular/year [get]
func (s *Store) GetPlaylistPopularByPastYear(c echo.Context) error {
	p, err := s.public.GetPlaylistPopularByPastYear(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public GetPlaylistByPastYear failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// GetPlaylistPopularByPastMonth returns a fake playlist with a list of popular videos of past month
//
// @Summary Provides a playlist of popular videos of past month
// @Description Provides a fake playlist, containing a list of popular videos
// @ID get-public-playlist-popular-by-past-month
// @Tags public-playlist
// @Produce json
// @Success 200 {object} public.Playlist
// @Router /v1/public/playlist/popular/month [get]
func (s *Store) GetPlaylistPopularByPastMonth(c echo.Context) error {
	p, err := s.public.GetPlaylistPopularByPastMonth(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public GetPlaylistByPastMonth failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// GetPlaylistRandom returns a fake playlist with a list of random videos
//
// @Summary Provides a playlist of random videos
// @Description Provides a fake playlist, containing a list of random videos
// @ID get-public-playlist-random
// @Tags public-playlist
// @Produce json
// @Success 200 {object} public.Playlist
// @Router /v1/public/playlist/random [get]
func (s *Store) GetPlaylistRandom(c echo.Context) error {
	p, err := s.public.GetPlaylistRandom(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("public GetPlaylistByRandom failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}
