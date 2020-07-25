package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

var config = echoMw.CORSConfig{
	AllowCredentials: true,
	Skipper:          echoMw.DefaultSkipper,
	AllowOrigins:     []string{"http://localhost:3000", "https://creator.ystv.co.uk", "http://comp.ystv.co.uk", "new.ystv.co.uk"},
	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}

// Init intialises web server middleware
func Init(e *echo.Echo) {
	e.Pre(echoMw.RemoveTrailingSlash())
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(config))
	e.Use(echoMw.Gzip())
}
