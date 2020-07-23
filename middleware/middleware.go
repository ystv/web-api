package middleware

import (
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

// Init intialises web server middleware
func Init(e *echo.Echo) {
	e.Pre(echoMw.RemoveTrailingSlash())
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORS())
	e.Use(echoMw.Gzip())
}
