package middleware

import (
	"github.com/ystv/web-api/utils",
	"github.com/labstack/echo/v4",
	echoMw "github.com/labstack/echo/v4/middleware",
)

// Init initialise database connection
func Init(e *echo.Echo) {
	utils.InitDB()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echoMw.Logger())
	e.Use(echoMw.recover())
}