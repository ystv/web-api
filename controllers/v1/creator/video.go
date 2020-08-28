package creator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/services/creator/types/video"
)

// VideoFind finds a video by ID
func (r *Repos) VideoFind(c echo.Context) error {
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

// VideoNew Handles creation of a video
func (r *Repos) VideoNew(c echo.Context) error {
	v := video.NewVideo{}
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

// VideoList Handles listing all creations
func (r *Repos) VideoList(c echo.Context) error {
	v, err := r.video.ListMeta(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("failed to list videos: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}

// VideosUser Handles retrieving a user's videos using their userid in their token.
func (r *Repos) VideosUser(c echo.Context) error {
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

// CalendarList Handles listing all videos from a calendar year/month
func (r *Repos) CalendarList(c echo.Context) error {
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
