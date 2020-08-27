package public

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

// Find handles converting a url path to either a video or series
func (r *Repos) Find(c echo.Context) error {
	raw := c.Request().URL
	rawSplit := strings.Split(raw.Path, "/")
	rawSplit = rawSplit[4:]
	rawJoined := strings.Join(rawSplit, "/")

	if len(rawJoined) == 0 {
		return c.String(http.StatusBadRequest, "Invalid URL, format [series]/[series]/[video]")
	}
	clean, err := url.Parse(rawJoined)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid URL")
	}
	b, err := r.public.Find(c.Request().Context(), clean.Path)
	if err != nil {
		log.Printf("PublicFind failed: %+v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, b)
}
