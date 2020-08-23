package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

var config = echoMw.CORSConfig{
	AllowCredentials: true,
	Skipper:          echoMw.DefaultSkipper,
	AllowOrigins: []string{
		"http://creator.ystv.co.uk",
		"https://creator.ystv.co.uk",
		"http://my.ystv.co.uk",
		"https://my.ystv.co.uk",
		"http://local.ystv.co.uk:3000",
		"http://localhost:3000"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}

// Init intialises web server middleware
func Init(e *echo.Echo) {
	e.Pre(echoMw.RemoveTrailingSlash())
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(config))
	e.Use(echoMw.GzipWithConfig(echoMw.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Path(), "swagger") {
				return true
			}
			return false
		},
	}))
}
