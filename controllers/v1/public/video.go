package public

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/public"
)

// Video handles a video item, providing info vfiles
func Video(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad video ID")
	}
	v, err := public.VideoFind(id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, v)
}

// ListVideos handles listing videos using an offset and page
func ListVideos(c echo.Context) error {
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad offset")
	}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad page")
	}
	v, err := public.VideoList(offset, page)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, v)
}

// URLToVideo handles converting the URL to a video
func URLToVideo(c echo.Context) error {
	public.URLToVideo(c.Request().URL)
	return c.NoContent(http.StatusOK)
}

// VideoBreadcrumb handles generating the breadcrumb of a video
func VideoBreadcrumb(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad video ID")
	}
	v, err := public.VideoBreadcrumb(id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, v)
}
