package creator

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	"github.com/ystv/web-api/services/creator/types/video"
	"github.com/ystv/web-api/utils"
)

// GetVideo finds a video by ID
//
// @Summary Get video by ID
// @Description Get a playlist including its children files.
// @ID get-creator-video
// @Tags creator-videos
// @Produce json
// @Param videoid path int true "Video ID"
// @Success 200 {object} video.Item
// @Router /v1/internal/creator/videos/{videoid} [get]
func (s *Store) GetVideo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid video ID")
	}

	v, err := s.video.GetItem(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("failed to get video item: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, v)
}

type NewVideoOutput struct {
	VideoID int `json:"id"`
}

// NewVideo Handles creation of a video
//
// @Summary New video
// @Description creates a new video, requires the file ID/name to find it in CDN.
// @ID new-creator-video
// @Tags creator-videos
// @Accept json
// @Param event body video.New true "NewVideo object"
// @Success 201 body int "Video ID"
// @Router /v1/internal/creator/videos [post]
func (s *Store) NewVideo(c echo.Context) error {
	var v video.New

	err := c.Bind(&v)
	if err != nil {
		err = fmt.Errorf("VideoCreate bind fail: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("VideoNew failed to get user ID: %w", err)
		return echo.NewHTTPError(status, err)
	}

	v.CreatedBy = claims.UserID

	videoID, err := s.video.NewItem(c.Request().Context(), v)
	if err != nil {
		err = fmt.Errorf("failed to create new video item: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, NewVideoOutput{VideoID: videoID})
}

// UpdateVideoMeta updates a video's metadata not files
//
// @Summary Update video meta
// @Description Updates a video metadata
// @ID update-creator-video-meta
// @Tags creator-videos
// @Accept json
// @Param event body video.Item true "VideoItem object"
// @Success 200 body int "Video ID"
// @Router /v1/internal/creator/video/meta [put]
func (s *Store) UpdateVideoMeta(c echo.Context) error {
	var v video.Meta

	err := c.Bind(&v)
	if err != nil {
		err = fmt.Errorf("failed to bind video object: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	t, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("failed to get token: %w", err)
		return echo.NewHTTPError(status, err)
	}

	parsed, err := strconv.ParseInt(strconv.Itoa(t.UserID), 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to parse user id: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	v.UpdatedByID = null.IntFrom(parsed)
	currentDateTime := time.Now()
	v.UpdatedAt = null.TimeFrom(currentDateTime)

	err = s.video.UpdateMeta(c.Request().Context(), v)
	if err != nil {
		err = fmt.Errorf("failed to update meta: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Store) DeleteVideo(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// ListVideos Handles listing all creations
//
// @Summary List all videos
// @Description Lists all videos, doesn't include files inside.
// @ID get-creator-videos-all
// @Tags creator-videos
// @Produce json
// @Success 200 {array} video.Meta
// @Router /v1/internal/creator/video [get]
func (s *Store) ListVideos(c echo.Context) error {
	v, err := s.video.ListMeta(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("failed to list videos: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}

// ListVideosByUser Handles retrieving a user's videos using their userid in their token.
//
// @Summary List all videos created by user ID
// @Description Lists all videos, doesn't include files inside.
// Use the createdBy user ID.
// @ID get-creator-videos-user
// @Tags creator-videos
// @Produce json
// @Success 200 {array} video.Meta
// @Router /v1/internal/creator/video/my [get]
func (s *Store) ListVideosByUser(c echo.Context) error {
	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("VideoNew failed to get user ID: %w", err)
		return echo.NewHTTPError(status, err)
	}

	v, err := s.video.ListMetaByUser(c.Request().Context(), claims.UserID)
	if err != nil {
		err = fmt.Errorf("failed to list videos: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}

// ListVideosByMonth Handles listing all videos from a calendar year/month
//
// @Summary List videos by month
// @Description Lists videos by month.
// @ID get-creator-videos-calendar
// @Tags creator-videos
// @Produce json
// @Param year path int true "year"
// @Param month path int true "month"
// @Success 200 {array} video.MetaCal
// @Router /v1/internal/creator/calendar/{year}/{month} [get]
func (s *Store) ListVideosByMonth(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Year incorrect, format /yyyy/mm")
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Month incorrect, format /yyyy/mm")
	}

	v, err := s.video.ListByCalendarMonth(c.Request().Context(), year, month)
	if err != nil {
		err = fmt.Errorf("failed to list by calendar month: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}

type searchInput struct {
	Query string `json:"query"`
}

// SearchVideo Handles listing appropriate videos from the relevant search
//
// @Summary Search videos
// @Description Search videos.
// @ID search-creator-videos
// @Tags creator-videos
// @Produce json
// @Param searchInput body searchInput true "Search Input object"
// @Success 200 {array} video.Meta
// @Router /v1/internal/creator/video/search [post]
func (s *Store) SearchVideo(c echo.Context) error {
	var input searchInput

	err := c.Bind(&input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	v, err := s.video.Search(c.Request().Context(), input.Query)
	if err != nil {
		err = fmt.Errorf("public Search failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(v))
}
