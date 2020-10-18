package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SeriesByID returns a series with it's immediate children with a SeriesID
//
// @Summary Provides a series
// @Description Returns a series object, including the children videos and series.
// @ID get-public-series
// @Tags public, series
// @Param seriesid path int true "Series ID"
// @Produce json
// @Success 200 {object} public.Series
// @Router /v1/public/series/{seriesid} [get]
func (r *Repos) SeriesByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad series ID")
	}
	s, err := r.public.GetSeries(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("Public SeriesByID failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, s)
}
