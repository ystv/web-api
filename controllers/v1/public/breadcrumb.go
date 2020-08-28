package public

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Find handles converting a url path to either a video or series
func (r *Repos) Find(c echo.Context) error {
	raw := c.Request().URL
	rawSplit := strings.Split(raw.Path, "/")
	rawSplit = rawSplit[4:]
	rawJoined := strings.Join(rawSplit, "/")

	if len(rawJoined) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL, format [series]/[series]/[video]")
	}
	clean, err := url.Parse(rawJoined)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL")
	}
	b, err := r.public.Find(c.Request().Context(), clean.Path)
	if err != nil {
		err = fmt.Errorf("Public Find failed: %w", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, b)
}

// VideoBreadcrumb handles generating the breadcrumb of a video
func (r *Repos) VideoBreadcrumb(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad video ID")
	}
	v, err := r.public.VideoBreadcrumb(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("Public VideoBreadcrumb failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}

// SeriesBreadcrumb returns the breadcrumb of a given series
func (r *Repos) SeriesBreadcrumb(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad series ID")
	}
	v, err := r.public.SeriesBreadcrumb(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("Public SeriesBreadcrumb failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}
