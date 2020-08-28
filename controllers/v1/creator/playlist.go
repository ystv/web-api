package creator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/services/creator/types/playlist"
	"gopkg.in/guregu/null.v4"
)

// PlaylistAll handles listing all playlist metadata's
func (r *Repos) PlaylistAll(c echo.Context) error {
	p, err := r.playlist.All(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("PlaylistAll failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PlaylistGet handles getting a single playlist and it's following videometa's
func (r *Repos) PlaylistGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
	}
	p, err := r.playlist.Get(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("Playlist get failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PlaylistNew handles creating a new playlist item
func (r *Repos) PlaylistNew(c echo.Context) error {
	p := playlist.Playlist{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate: failed to bind json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// TODO sort out user ID
	res, err := r.playlist.New(c.Request().Context(), p)
	if err != nil {
		err = fmt.Errorf("PlaylistNew failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// PlaylistUpdate handles updating a playlist
func (r *Repos) PlaylistUpdate(c echo.Context) error {
	p := playlist.Playlist{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate: failed to bind json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	claims, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	p.UpdatedBy = null.IntFrom(int64(claims.UserID))
	var videoIDs []int
	for _, v := range p.Videos {
		videoIDs = append(videoIDs, v.ID)
	}
	err = r.playlist.Update(c.Request().Context(), p.Meta, videoIDs)
	if err != nil {
		err = fmt.Errorf("PlaylistUpdate: failed to update playlist: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
