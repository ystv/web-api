package creator

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// SeriesListAll handles listing every series and their depth
func (r *Repos) SeriesListAll(c echo.Context) error {
	s, err := r.series.All(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, s)
	}
	return c.JSON(http.StatusOK, s)
}
