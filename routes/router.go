package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	v1creator "github.com/ystv/web-api/controllers/v1/creator"
	v1playlist "github.com/ystv/web-api/controllers/v1/playlist"
	v1stream "github.com/ystv/web-api/controllers/v1/stream"
	v1tables "github.com/ystv/web-api/controllers/v1/tables"
	v1video "github.com/ystv/web-api/controllers/v1/video"
	_ "github.com/ystv/web-api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ystv/web-api/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// Init initialise routes
func Init() *echo.Echo {
	e := echo.New()

	debug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		panic(err)
	}
	e.Debug = debug

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
			videoBoxes := tables.Group("/video_boxes")
			{
				videoBoxes.POST("", v1tables.VideoBoxCreate)
				videoBoxes.GET("", v1tables.VideoBoxList)
				videoBoxes.GET("/:id", v1tables.VideoBoxFind)
				videoBoxes.PUT("/:id", v1tables.VideoBoxUpdate)
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
		internal := apiV1.Group("/internal")
		{
			creator := internal.Group("/creator")
			{
				creation := creator.Group("/:id")
				{
					creation.POST("/meta", v1creator.CreationMetaCreate)
					creation.POST("/video", v1creator.CreationFileUpload)
					creator.GET("", v1creator.CreationFind)
				}
				creator.POST("", v1creator.CreationCreate)
				creator.GET("", v1creator.CreationList)
			}
		}
		public := apiV1.Group("/public")
		{
			videos := public.Group("/videos/:id")
			{
				videos.GET("", v1video.Info)
				videos.GET("/full", v1video.Full)
			}
			playlists := public.Group("/playlist/:id")
			{
				playlists.GET("", v1playlist.Info)
			}
		}
		stream := apiV1.Group("/stream")
		{
			stream.POST("/auth", v1stream.CheckAuth)
		}

	}
	e.GET("/", func(c echo.Context) error {
		text := `                                                                                
                                                              @@@@@             
                                                                     @@@@       
                                                                         @@@    
                                               @@@@                        @@@@ 
                                               @@@@                          @@@
        .    @@@@@         @@@@   @@@@@@     @@@@@@@@@  @@@@        @@@@     @@@
     @%       @@@@@       @@@@  @@@@@@@@@@   @@@@@@@@@   @@@@      @@@@       @@
   @@           @@@@    @@@@@   @@@@           @@@@       @@@@    @@@@       @@ 
  @@             @@@@  @@@@@     @@@@@@@@      @@@@        @@@@  @@@@       @@  
 @@               @@@@@@@@           @@@@@@    @@@@         @@@@@@@@      @@    
 @@@               @@@@@@       @@@    @@@@    @@@@          @@@@@@     @       
 @@@                @@@@        @@@@@@@@@@     @@@@           @@@@              
  (@@@             @@@@                                                         
     @@@         @@@@@                                                          
        @@@@    @@@@@                                                           
              @@@@@                                                             
`
		return c.String(http.StatusOK, text)
	})
	return e
}
