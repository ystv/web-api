package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Video handles a video item, providing info vfiles
func (r *Repos) Video(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad video ID")
	}
	v, err := r.public.GetVideo(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("Public GetVideo failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}

// ListVideos handles listing videos using an offset and page
func (r *Repos) ListVideos(c echo.Context) error {
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad offset")
	}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad page")
	}
	v, err := r.public.ListVideo(c.Request().Context(), offset, page)
	if err != nil {
		err = fmt.Errorf("Public ListVideos failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}
