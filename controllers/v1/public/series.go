package public

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/public"
)

// SeriesByID returns a series with it's immediate children with a SeriesID
func SeriesByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad video ID")
	}
	v, err := public.SeriesAndChildren(id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, v)
}
