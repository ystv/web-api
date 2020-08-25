package creator

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/types/playlist"
)

// PlaylistAll handles listing all playlist metadata's
func (r *Repos) PlaylistAll(c echo.Context) error {
	p, err := r.playlist.All(c.Request().Context())
	if err != nil {
		log.Printf("Playlist all failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, p)
	}
	return c.JSON(http.StatusOK, p)
}

// PlaylistGet handles getting a single playlist and it's following videometa's
func (r *Repos) PlaylistGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	p, err := r.playlist.Get(c.Request().Context(), id)
	if err != nil {
		log.Printf("Playlist get failed: %+v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PlaylistNew handles creating a new playlist item
func (r *Repos) PlaylistNew(c echo.Context) error {
	p := playlist.Playlist{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	res, err := r.playlist.New(c.Request().Context(), p)
	if err != nil {
		log.Printf("Playlist new failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}
