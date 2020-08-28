package creator

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ListSeries handles listing every series and their depth
func (r *Repos) ListSeries(c echo.Context) error {
	s, err := r.series.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, s)
	}
	return c.JSON(http.StatusOK, s)
}

// GetSeries finds a video by ID
func (r *Repos) GetSeries(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	v, err := r.series.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, v)
}
