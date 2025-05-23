package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils"
)

// GetVideo handles a video item, providing info
//
// @Summary Provides a video item
// @Description Returns a video item. Including the video files.
// @ID get-public-video
// @Tags public-video
// @Param videoid path int true "Video ID"
// @Produce json
// @Success 200 {object} public.VideoItem
// @Router /v1/public/video/{videoid} [get]
func (s *Store) GetVideo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad video ID")
	}

	v, err := s.public.GetVideo(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("public GetVideo failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, v)
}

// ListVideos handles listing videos using an offset and page
//
// @Summary Provides a list of videos
// @Description List of video meta's in order of broadcast date.
// @ID get-public-videos
// @Tags public-video
// @Param offset path int true "Offset"
// @Param page path int true "Page"
// @Produce json
// @Success 200 {array} public.VideoMeta
// @Router /v1/public/videos/{offset}/{page} [get]
func (s *Store) ListVideos(c echo.Context) error {
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad offset")
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad page")
	}

	v, err := s.public.ListVideos(c.Request().Context(), offset, page)
	if err != nil {
		err = fmt.Errorf("public ListVideos failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}
