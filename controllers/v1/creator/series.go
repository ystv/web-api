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
func (r *Repos) ListSeries(c echo.Context) error {
	s, err := r.series.List(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("failed to list series: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(s))
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
func (r *Repos) GetSeries(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Number pls")
	}

	s, err := r.series.GetSeries(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("failed to get series: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, s)
}

// UpdateSeries handles updating a series
// @Summary UpdatePlaylist series
// @Description UpdatePlaylist a series, video ID's required otherwise it will remove all videos.
// @ID update-creator-series
// @Tags creator-series
// @Accept json
// @Param quote body series.NewPlaylist true "Series object"
// @Success 200
// @Router /v1/internal/creator/series [put]
func (r *Repos) UpdateSeries(c echo.Context) error {
	var s series.Series

	err := c.Bind(&s)
	if err != nil {
		err = fmt.Errorf("SeriesUpdate: failed to bind json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("SeriesUpdate failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	_ = claims

	//s.UpdatedBy = &claims.UserID
	//
	//var videoIDs []int
	//for _, v := range s.Videos {
	//	videoIDs = append(videoIDs, v.ID)
	//}
	//
	//err = r.series.Update(c.Request().Context(), s.Meta, videoIDs)
	//if err != nil {
	//	err = fmt.Errorf("SeriesUpdate: failed to update series: %w", err)
	//	return echo.NewHTTPError(http.StatusInternalServerError, err)
	//}

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
func (r *Repos) DeleteSeries(c echo.Context) error {
	seriesID, err := strconv.Atoi(c.Param("seriesid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	//err = r.series.DeleteSeries(c.Request().Context(), seriesID)
	//if err != nil {
	//	err = fmt.Errorf("DeleteSeries failed: %w", err)
	//	return c.JSON(http.StatusInternalServerError, err)
	//}

	_ = seriesID

	return c.NoContent(http.StatusOK)
}
