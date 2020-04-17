package routes

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	v1tables "github.com/ystv/web-api/controllers/v1/tables"
	_ "github.com/ystv/web-api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ystv/web-api/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// Init initialise routes
func Init() *echo.Echo {
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
		tables := apiV1.Group("/tables")
		{
			videos := tables.Group("/videos")
			{
				videos.POST("", v1tables.VideoCreate)
				videos.GET("", v1tables.VideoList)
				videos.GET("/:id", v1tables.VideoFind)
				videos.PUT("/:id", v1tables.VideoUpdate)
				videos.DELETE("/:id", v1tables.VideoDelete)
			}
			quotes := tables.Group("/quotes")
			{
				quotes.POST("", v1tables.QuoteCreate)
				quotes.GET("", v1tables.QuoteList)
				quotes.GET("/:id", v1tables.QuoteFind)
				quotes.PUT("/:id", v1tables.QuoteUpdate)
				quotes.DELETE("/:id", v1tables.QuoteDelete)
			}
		}

	}
	e.GET("/", func(c echo.Context) error {
		content, err := ioutil.ReadFile("logo.txt")
		if err != nil {
			panic(err)
		}
		text := string(content)
		return c.String(http.StatusOK, text)
	})
	return e
}
