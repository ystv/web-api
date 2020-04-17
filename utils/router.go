package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	v1 "github.com/ystv/web-api/controllers/v1"
	"github.com/ystv/web-api/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitRoutes() *echo.Echo {
	e := echo.New()

	e.Debug = true

	middleware.Init(e)

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/routes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, e.Routes())
	})
	// ping
	e.GET("/ping", func(c echo.Context) error {
		resp := map[string]time.Time{"pong": time.Now()}
		return c.JSON(http.StatusOK, resp)
	})

	apiV1 := e.Group("v1")
	{
		videos := apiV1.Group("/videos")
		{
			videos.POST("", v1.VideoCreate)
			videos.GET("", v1.VideoList)
			videos.PUT("/:id", v1.VideoUpdate)
			videos.DELETE("/:id", v1.VideoDelete)
		}
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "YSTV API")
	})
	return e
}