package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SeriesByID returns a series with it's immediate children with a SeriesID
func (r *Repos) SeriesByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad series ID")
	}
	v, err := r.public.GetSeries(c.Request().Context(), id)
	if err != nil {
		err = fmt.Errorf("Public SeriesByID failed : %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, v)
}
