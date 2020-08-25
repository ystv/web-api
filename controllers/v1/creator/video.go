package creator

import (
	"log"
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
		c.String(http.StatusBadRequest, "Number pls")
	}
	v, err := r.video.GetItem(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
}

// VideoNew Handles creation of a video
func (r *Repos) VideoNew(c echo.Context) error {
	v := video.NewVideo{}
	err := c.Bind(&v)
	if err != nil {
		log.Printf("VideoCreate bind fail: %+v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	claims, err := people.GetToken(c)
	if err != nil {
		log.Printf("VideoNew failed to get user ID: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	v.CreatedBy = claims.UserID
	err = r.video.NewItem(c.Request().Context(), &v)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, "Creation created")
}

// VideoList Handles listing all creations
func (r *Repos) VideoList(c echo.Context) error {
	creations, err := r.video.ListMeta(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, creations)
}

// VideosUser Handles retrieving a user's videos using their userid in their token.
func (r *Repos) VideosUser(c echo.Context) error {
	claims, err := people.GetToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	v, err := r.video.ListMetaByUser(c.Request().Context(), claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}

// CalendarList Handles listing all videos from a calendar year/month
func (r *Repos) CalendarList(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Year incorrect, format /yyyy/mm")
	}
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Month incorrect, format /yyyy/mm")
	}
	v, err := r.video.ListByCalendarMonth(c.Request().Context(), year, month)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
}
