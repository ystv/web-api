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

// SeriesByYear returns a virtual series containing series / videos made in that year
//
// @Summary Series of a year
// @Description Returns a series array, virtual series that contains child series / videos
// @Description that were made in that year.
// @ID get-public-series-year
// @Tags public, series
// @Param year path int true "Year"
// @Produce json
// @Success 200 {array} public.Series
// @Router /v1/public/series/yearly/{year} [get]
func (r *Repos) SeriesByYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}
	s, err := r.public.SeriesByYear(c.Request().Context(), year)
	if err != nil {
		err = fmt.Errorf("Public ListSeriesByYear failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, s)
}
