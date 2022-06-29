package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

// New intialises web server middleware
func New(e *echo.Echo, domainName string) {
	config := echoMw.CORSConfig{
		AllowCredentials: true,
		Skipper:          echoMw.DefaultSkipper,
		AllowOrigins: []string{
			"http://creator." + domainName,
			"https://creator." + domainName,
			"http://my." + domainName,
			"https://my." + domainName,
			"http://local." + domainName + ":3000",
			"https://local." + domainName + ":3000",
			"http://local." + domainName + ":8080",
			"https://local." + domainName + ":8080",
      "ystv-development.localhost:3000",
			"http://" + strings.Join(strings.Split(domainName, ".")[1:], "."),
			"https://" + strings.Join(strings.Split(domainName, ".")[1:], "."),
			"http://" + domainName,
			"https://" + domainName},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials, echo.HeaderAccessControlAllowOrigin, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Pre(echoMw.RemoveTrailingSlash())
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(config))
	e.Use(echoMw.GzipWithConfig(echoMw.GzipConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			if path == "/metrics" || strings.Contains(path, "swagger") {
				return true
			}
			return false
		},
	}))
	// TODO secure this
	// /metrics, view using curl
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
}
