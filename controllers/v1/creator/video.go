package creator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/services/creator/types/video"
)

// GetVideo finds a video by ID
//
// @Summary Get video by ID
// @Description Get a playlist including it's children files.
// @ID get-creator-video
// @Tags creator-videos
// @Produce json
// @Param videoid path int true "Video ID"
// @Success 200 {object} video.Item
// @Router /v1/internal/creator/videos/{videoid} [get]
func (r *Repos) GetVideo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "Invalid video ID")
	}
	v, err := r.video.GetItem(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("failed to get video item: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
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
	v := video.New{}
	err := c.Bind(&v)
	if err != nil {
		err = fmt.Errorf("VideoCreate bind fail: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	claims, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("VideoNew failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	v.CreatedBy = claims.UserID
	err = r.video.NewItem(c.Request().Context(), &v)
	if err != nil {
		err = fmt.Errorf("failed to create new video item: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// TODO return created video ID
	return c.String(http.StatusCreated, "Creation created")
}

// UpdateVideo updates a video's metadata not files
func (r *Repos) UpdateVideo(c echo.Context) error {
	v := video.Item{}
	err := c.Bind(&v)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}
	return c.NoContent(http.StatusOK)
}

func (r *Repos) DeleteVideo(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// VideoList Handles listing all creations
//
// @Summary List all videos
// @Description Lists all videos, doesn't include files inside.
// @ID get-creator-videos-all
// @Tags creator-videos
// @Produce json
// @Success 200 {array} video.Meta
// @Router /v1/internal/creator/videos [get]
func (r *Repos) VideoList(c echo.Context) error {
	v, err := r.video.ListMeta(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("failed to list videos: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}

// ListVideosByUser Handles retrieving a user's videos using their userid in their token.
//
// @Summary List all videos created by user ID
// @Description Lists all videos, doesn't include files inside. Uses the createdBy user ID.
// @ID get-creator-videos-user
// @Tags creator-videos
// @Produce json
// @Success 200 {array} video.Meta
// @Router /v1/internal/creator/videos/my [get]
func (r *Repos) ListVideosByUser(c echo.Context) error {
	claims, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("VideoNew failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	v, err := r.video.ListMetaByUser(c.Request().Context(), claims.UserID)
	if err != nil {
		err = fmt.Errorf("failed to list videos: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
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
	return c.JSON(http.StatusOK, v)
}
