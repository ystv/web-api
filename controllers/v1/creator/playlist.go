package creator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/creator/types/playlist"
	"github.com/ystv/web-api/utils"
)

// ListPlaylists handles listing all playlist metadata
// @Summary List all playlists
// @Description Lists all playlists, doesn't include videos inside.
// @ID get-creator-playlists-all
// @Tags creator-playlists
// @Produce json
// @Success 200 {array} playlist.Playlist
// @Router /v1/internal/creator/playlist [get]
func (r *Repos) ListPlaylists(c echo.Context) error {
	p, err := r.playlist.ListPlaylists(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("PlaylistAll failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// GetPlaylist handles getting a single playlist, and it's following videometa's
// @Summary Get playlist by ID
// @Description Get a playlist including its children videos.
// @ID get-creator-playlist
// @Tags creator-playlists
// @Produce json
// @Param playlistid path int true "Playlist ID"
// @Success 200 {object} playlist.Playlist
// @Router /v1/internal/creator/playlist/{playlistid} [get]
func (r *Repos) GetPlaylist(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
	}

	p, err := r.playlist.GetPlaylist(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("playlist get failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, p)
}

// NewPlaylist handles creating a new playlist item
// @Summary New playlist
// @Description creates a new playlist with optional video ID's.
// @ID new-creator-playlist
// @Tags creator-playlists
// @Accept json
// @Param event body playlist.Playlist true "Playlist object"
// @Success 201 body int "Playlist ID"
// @Router /v1/internal/creator/playlist [post]
func (r *Repos) NewPlaylist(c echo.Context) error {
	var p playlist.New

	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("failed to bind json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	p.CreatedBy = claims.UserID

	res, err := r.playlist.NewPlaylist(c.Request().Context(), p)
	if err != nil {
		err = fmt.Errorf("PlaylistNew failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdatePlaylist handles updating a playlist
// @Summary Update playlist
// @Description Update a playlist, video ID's required otherwise it will remove all videos.
// @ID update-creator-playlist
// @Tags creator-playlists
// @Accept json
// @Param quote body playlist.New true "Playlist object"
// @Success 200
// @Router /v1/internal/creator/playlist [put]
func (r *Repos) UpdatePlaylist(c echo.Context) error {
	var p playlist.Playlist

	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate: failed to bind json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	p.UpdatedBy = &claims.UserID

	videoIDs := make([]int, 0)
	for _, v := range p.Videos {
		videoIDs = append(videoIDs, v.ID)
	}

	err = r.playlist.UpdatePlaylist(c.Request().Context(), p.Meta, videoIDs)
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate: failed to update playlist: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// DeletePlaylist handles deleting playlist
// @Summary Delete a playlist
// @Description Delete a playlist
// @ID delete-creator-playlist
// @Tags creator-playlist
// @Param playlistid path int true "Series ID"
// @Success 200
// @Router /v1/internal/creator/playlist/{playlistid} [delete]
func (r *Repos) DeletePlaylist(c echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlistid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	err = r.playlist.DeletePlaylist(c.Request().Context(), playlistID)
	if err != nil {
		err = fmt.Errorf("DeletePlaylist failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
