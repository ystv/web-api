package middleware

import (
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"github.com/ystv/web-api/utils"
)

// Init initialse database connection
func Init(e *echo.Echo) {
	utils.InitDB()
	utils.InitCDN()
	utils.InitMessaging()
	// utils.InitAuth()

	e.Pre(echoMw.RemoveTrailingSlash())
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORS())
	e.Use(echoMw.Gzip())
}
