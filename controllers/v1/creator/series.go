package creator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/utils"
)

// ListSeries handles listing every series and their depth
// @Summary List all series
// @Description Lists all series, doesn't include videos inside.
// @ID get-creator-series-all
// @Tags creator-series
// @Produce json
// @Success 200 {array} series.Meta
// @Router /v1/internal/creator/series [get]
func (s *Store) ListSeries(c echo.Context) error {
	series1, err := s.series.List(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("failed to list series: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(series1))
}

// GetSeries finds a video by ID
// @Summary Get series by ID
// @Description Get a series including its children videos.
// @ID get-creator-series
// @Tags creator-series
// @Produce json
// @Param seriesid path int true "Series ID"
// @Success 200 {object} series.Series
// @Router /v1/internal/creator/series/{seriesid} [get]
func (s *Store) GetSeries(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Number pls")
	}

	series1, err := s.series.GetSeries(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("failed to get series: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, series1)
}

// UpdateSeries handles updating a series
// @Summary UpdatePlaylist series
// @Description UpdatePlaylist a series, video ID's required otherwise it will remove all videos.
// @ID update-creator-series
// @Tags creator-series
// @Accept json
// @Param quote body series.Series true "Series object"
// @Success 200
// @Router /v1/internal/creator/series [put]
func (s *Store) UpdateSeries(c echo.Context) error {
	var series1 series.Series

	err := c.Bind(&series1)
	if err != nil {
		err = fmt.Errorf("SeriesUpdate: failed to bind json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("SeriesUpdate failed to get user ID: %w", err)
		return echo.NewHTTPError(status, err)
	}

	_ = claims

	// series1.UpdatedBy = &claims.UserID
	//
	// var videoIDs []int
	// for _, v := range series1.Videos {
	//	videoIDs = append(videoIDs, v.ID)
	// }
	//
	// err = s.series.Update(c.Request().Context(), series1.Meta, videoIDs)
	// if err != nil {
	//	err = fmt.Errorf("SeriesUpdate: failed to update series: %w", err)
	//	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	return c.NoContent(http.StatusOK)
}

// DeleteSeries handles deleting series
// @Summary Delete a series
// @Description Delete a series
// @ID delete-creator-series
// @Tags creator-series
// @Param seriesid path int true "Series ID"
// @Success 200
// @Router /v1/internal/creator/series/{seriesid} [delete]
func (s *Store) DeleteSeries(c echo.Context) error {
	seriesID, err := strconv.Atoi(c.Param("seriesid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	// err = s.series.DeleteSeries(c.Request().Context(), seriesID)
	// if err != nil {
	//	err = fmt.Errorf("DeleteSeries failed: %w", err)
	//	return c.JSON(http.StatusInternalServerError, err)
	// }

	_ = seriesID

	return c.NoContent(http.StatusOK)
}
