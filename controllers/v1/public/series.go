package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetSeriesByID returns a series with its immediate children with a SeriesID
//
// @Summary Provides a series
// @Description Returns a series object, including the children videos and series.
// @ID get-public-series
// @Tags public-series
// @Param seriesid path int true "Series ID"
// @Produce json
// @Success 200 {object} public.Series
// @Router /v1/public/series/{seriesid} [get]
func (s *Store) GetSeriesByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad series ID")
	}
	series, err := s.public.GetSeries(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("public SeriesByID failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, series)
}

// GetSeriesByYear returns a virtual series containing series / videos made in that year
//
// @Summary Series of a year
// @Description Returns a series array, virtual series that contains child series / videos
// @Description that were made in that year.
// @ID get-public-series-year
// @Tags public-series
// @Param year path int true "Year"
// @Produce json
// @Success 200 {array} public.Series
// @Router /v1/public/series/yearly/{year} [get]
func (s *Store) GetSeriesByYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad year")
	}

	series, err := s.public.GetSeriesByYear(c.Request().Context(), year)
	if err != nil {
		err = fmt.Errorf("public ListSeriesByYear failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, series)
}

type SearchInput struct {
	Query string `json:"query"`
}

// Search returns a virtual series that contains relevant videos and series
//
// @Summary Search the VOD library
// @Description Returns a virtual series that contains relevant videos and series
// @ID search-vod
// @Tags public-series
// @Param searchInput body SearchInput true "Search Input object"
// @Produce json
// @Success 200 {array} public.Series
// @Router /v1/public/search [post]
func (s *Store) Search(c echo.Context) error {
	var searchInput SearchInput

	err := c.Bind(&searchInput)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	series, err := s.public.Search(c.Request().Context(), searchInput.Query)
	if err != nil {
		err = fmt.Errorf("public Search failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, series)
}
