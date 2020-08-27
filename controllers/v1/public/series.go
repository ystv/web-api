package public

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SeriesByID returns a series with it's immediate children with a SeriesID
func (r *Repos) SeriesByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad video ID")
	}
	v, err := r.public.GetSeries(c.Request().Context(), id)
	if err != nil {
		log.Printf("Public SeriesByID failed : %+v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, v)
}

// SeriesBreadcrumb returns the breadcrumb of a given series
func (r *Repos) SeriesBreadcrumb(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad video ID")
	}
	v, err := r.public.SeriesBreadcrumb(c.Request().Context(), id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, v)
}
