package playlist

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Info Handles playlist info i.e. Title, thumbnail, desc
func Info(c echo.Context) error {
	return c.String(http.StatusOK, "Playlist info")
}
