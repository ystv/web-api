package video

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Info Handles video info i.e. Title, thumbnail, desc
func Info(c echo.Context) error {
	return c.String(http.StatusOK, "Video info")
}

// Full Handles all video information i.e. files
func Full(c echo.Context) error {
	return c.String(http.StatusOK, "Video full")
}
