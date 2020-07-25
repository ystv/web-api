package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	v1creator "github.com/ystv/web-api/controllers/v1/creator"
	v1public "github.com/ystv/web-api/controllers/v1/public"
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
func Init() *echo.Echo {
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
					//videos.GET("", v1creator.VideoList)
					//videos.POST("", v1creator.VideoCreate)
					videoItem := videos.Group("/:id")
					{
						videoItem.GET("", v1creator.VideoFind)
						videoItem.PUT("", v1creator.CreationCreate)
					}
				}
				playlists := creator.Group("/playlists")
				{
					playlists.GET("", v1creator.CreationCreate)
					playlists.POST("", v1creator.CreationCreate)
				}
				creation := creator.Group("/:id")
				{
					creation.POST("/meta", v1creator.CreationMetaCreate)
					creation.POST("/video", v1creator.CreationFileUpload)
					creation.GET("", v1creator.VideoFind)
				}
				creator.POST("", v1creator.CreationCreate)
				creator.GET("", v1creator.CreationList)
				creator.GET("/calendar/:year/:month", v1creator.CalendarList)
				creator.GET("/stats", v1creator.Stats)
			}
		}
		public := apiV1.Group("/public")
		{
			public.GET("/videos/:offset/:page", v1public.ListVideos)
			public.GET("/video/:id", v1public.Video)
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
