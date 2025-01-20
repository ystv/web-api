package public

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/public"
	"github.com/ystv/web-api/utils"
)

// Find handles converting an url path to either a video or series
//
// @Summary Converts a VOD url to either a series or video
// @Description Allows us to remain backwards compatible with the existing URLs
// @ID get-public-breadcrumb-find
// @Tags public-breadcrumb
// @Param url path string true "URL Path"
// @Produce json
// @Success 200 {object} public.BreadcrumbItem
// @Router /v1/public/find/{url} [get]
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
		if errors.Is(err, public.ErrVideoNotFound) ||
			errors.Is(err, public.ErrSeriesNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = fmt.Errorf("public find failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, b)
}

// VideoBreadcrumb handles generating the breadcrumb of a video
//
// @Summary Provides a breadcrumb for a video
// @Description Returns a path of series to a video
// @ID get-public-breadcrumb-video
// @Tags public-breadcrumb
// @Param videoid path int true "Video ID"
// @Produce json
// @Success 200 {object} public.Breadcrumb
// @Router /v1/public/video/{videoid}/breadcrumb [get]
func (r *Repos) VideoBreadcrumb(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad video ID")
	}

	v, err := r.public.GetVideoBreadcrumb(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, public.ErrVideoNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = fmt.Errorf("public video breadcrumb failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}

// GetSeriesBreadcrumb returns the breadcrumb of a given series
//
// @Summary Provides a breadcrumb for a series
// @Description Returns a path of series to a series
// @ID get-public-breadcrumb-series
// @Tags public-breadcrumb
// @Param seriesid path int true "Series ID"
// @Produce json
// @Success 200 {object} public.Breadcrumb
// @Router /v1/public/series/{seriesid}/breadcrumb [get]
func (r *Repos) GetSeriesBreadcrumb(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad series ID")
	}

	v, err := r.public.GetSeriesBreadcrumb(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, public.ErrSeriesNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = fmt.Errorf("public series breadcrumb failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}
