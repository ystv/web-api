package routes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	creatorV1 "github.com/ystv/web-api/controllers/v1/creator"
	publicV1 "github.com/ystv/web-api/controllers/v1/public"
	v1stream "github.com/ystv/web-api/controllers/v1/stream"
	_ "github.com/ystv/web-api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ystv/web-api/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// JWTClaims represents an identifiable JWT
type JWTClaims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// Init initialise routes
func Init(version, commit string) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	debug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		panic(err)
	}
	e.Debug = debug

	middleware.Init(e)
	config := echoMw.JWTConfig{
		Claims:      &JWTClaims{},
		TokenLookup: "cookie:token",
		SigningKey:  []byte(os.Getenv("signing_key")),
	}

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
		internal := apiV1.Group("/internal")
		if !debug {
			internal.Use(echoMw.JWTWithConfig(config))
		}
		{
			creator := internal.Group("/creator")
			{
				videos := creator.Group("/videos")
				{
					//videos.GET("", creatorV1.VideoList)
					//videos.POST("", creatorV1.VideoCreate)
					videoItem := videos.Group("/:id")
					{
						videoItem.GET("", creatorV1.VideoFind)
						videoItem.PUT("", creatorV1.CreationCreate)
					}
				}
				playlists := creator.Group("/playlists")
				{
					playlists.GET("", notImplemented)
					playlists.POST("", notImplemented)
				}
				creation := creator.Group("/:id")
				{
					creation.POST("/meta", creatorV1.CreationMetaCreate)
					creation.POST("/video", creatorV1.CreationFileUpload)
					creation.GET("", creatorV1.VideoFind)
				}
				creator.POST("", notImplemented)
				creator.GET("", creatorV1.CreationList)
				creator.GET("/calendar/:year/:month", creatorV1.CalendarList)
				creator.GET("/stats", creatorV1.Stats)
			}
		}
		public := apiV1.Group("/public")
		{
			public.GET("/videos/:offset/:page", publicV1.ListVideos)
			public.GET("/video/:id", publicV1.Video)
			public.GET("/video/:id/breadcrumb", publicV1.VideoBreadcrumb)
			public.GET("/video/by_url", publicV1.URLToVideo)
			public.GET("/series/:id", publicV1.SeriesByID)
			public.GET("/series/:id/breadcrumb", publicV1.SeriesBreadcrumb)
			public.GET("/teams", publicV1.ListTeams)
		}
		stream := apiV1.Group("/stream")
		{
			stream.POST("/auth", v1stream.CheckAuth)
		}

	}
	e.GET("/", func(c echo.Context) error {
		text := fmt.Sprintf(`                                                                                
                                                              @@@@@             
                                                                     @@@@       
                                                                         @@@    
                                               @@@@                        @@@@ 
                                               @@@@                          @@@
        .    @@@@@         @@@@   @@@@@@     @@@@@@@@@  @@@@        @@@@     @@@
     @@       @@@@@       @@@@  @@@@@@@@@@   @@@@@@@@@   @@@@      @@@@       @@
   @@           @@@@    @@@@@   @@@@           @@@@       @@@@    @@@@       @@ 
  @@             @@@@  @@@@@     @@@@@@@@      @@@@        @@@@  @@@@       @@  
 @@               @@@@@@@@           @@@@@@    @@@@         @@@@@@@@      @@    
 @@@               @@@@@@       @@@    @@@@    @@@@          @@@@@@     @       
 @@@                @@@@        @@@@@@@@@@     @@@@           @@@@              
  (@@@             @@@@                                                         
     @@@         @@@@@          web-api                                         
        @@@@    @@@@@           Version: %s                                     
              @@@@@             Commit ID: %s                                                
`, version, commit)
		return c.String(http.StatusOK, text)
	})
	return e
}

func notImplemented(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
