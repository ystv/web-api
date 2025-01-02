package creator

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

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
func (r *Repos) GetVideo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid video ID")
	}

	v, err := r.video.GetItem(c.Request().Context(), id)
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
func (r *Repos) NewVideo(c echo.Context) error {
	var v video.New

	err := c.Bind(&v)
	if err != nil {
		err = fmt.Errorf("VideoCreate bind fail: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("VideoNew failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	v.CreatedBy = claims.UserID

	videoID, err := r.video.NewItem(c.Request().Context(), v)
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
func (r *Repos) UpdateVideoMeta(c echo.Context) error {
	var v video.Meta

	err := c.Bind(&v)
	if err != nil {
		err = fmt.Errorf("failed to bind video object: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	t, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	v.UpdatedByID = &t.UserID
	currentDateTime := time.Now()
	v.UpdatedAt = &currentDateTime

	err = r.video.UpdateMeta(c.Request().Context(), v)
	if err != nil {
		err = fmt.Errorf("failed to update meta: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (r *Repos) DeleteVideo(c echo.Context) error {
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
func (r *Repos) ListVideos(c echo.Context) error {
	v, err := r.video.ListMeta(c.Request().Context())
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
func (r *Repos) ListVideosByUser(c echo.Context) error {
	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("VideoNew failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	v, err := r.video.ListMetaByUser(c.Request().Context(), claims.UserID)
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
func (r *Repos) ListVideosByMonth(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Year incorrect, format /yyyy/mm")
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Month incorrect, format /yyyy/mm")
	}

	v, err := r.video.ListByCalendarMonth(c.Request().Context(), year, month)
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
func (r *Repos) SearchVideo(c echo.Context) error {
	var input searchInput

	err := c.Bind(&input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	s, err := r.video.Search(c.Request().Context(), input.Query)
	if err != nil {
		err = fmt.Errorf("public Search failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(s))
}
