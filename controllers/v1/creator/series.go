package creator

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator/series"
)

// SeriesListAll handles listing every series and their depth
func SeriesListAll(c echo.Context) error {
	s, err := series.All()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, s)
	}
	return c.JSON(http.StatusOK, s)
}
