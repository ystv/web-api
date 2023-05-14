package creator

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
		return c.JSON(http.StatusInternalServerError, s)
	}
	return c.JSON(http.StatusOK, s)
}

// GetSeries finds a video by ID
// @Summary Get series by ID
// @Description Get a series including it's children videos.
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
	s, err := r.series.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, s)
}
