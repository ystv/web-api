package creator

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/series"
)

// VideoMetaCreate Handes uploading meta data for a creation
func VideoMetaCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Meta created")
}

// FileUpload Handles uploading a file
func FileUpload(c echo.Context) error {
	creator.CreateBucket("pending", "ystv-wales-1")
	url, err := creator.GenerateUploadURL("pending", c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, url)
}

// VideoFind finds a video by ID
func VideoFind(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	v, err := creator.VideoItemFind(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
}

// VideoCreate Handles creation of a creation lol
func VideoCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Creation created")
}

// VideoList Handles listing all creations
func VideoList(c echo.Context) error {
	creations, err := creator.VideoMetaList(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, creations)
}

// CalendarList Handles listing all videos from a calendar year/month
func CalendarList(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Year incorrect, format /yyyy/mm")
	}
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Month incorrect, format /yyyy/mm")
	}
	v, err := creator.CalendarList(context.Background(), year, month)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
}

// Stats handles sending general stats about the video library
func Stats(c echo.Context) error {
	s, err := creator.Stats(context.Background())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, s)
}

// SeriesListAll handles listing every series and their depth
func SeriesListAll(c echo.Context) error {
	s, err := series.All()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, s)
	}
	return c.JSON(http.StatusOK, s)
}
