package routes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	clapperPackage "github.com/ystv/web-api/controllers/v1/clapper"
	creatorPackage "github.com/ystv/web-api/controllers/v1/creator"
	encoderV1 "github.com/ystv/web-api/controllers/v1/encoder"
	miscPackage "github.com/ystv/web-api/controllers/v1/misc"
	peoplePackage "github.com/ystv/web-api/controllers/v1/people"
	publicPackage "github.com/ystv/web-api/controllers/v1/public"
	streamV1 "github.com/ystv/web-api/controllers/v1/stream"
	_ "github.com/ystv/web-api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ystv/web-api/middleware"
	"github.com/ystv/web-api/utils"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// TODO standarise on function names

// Init initialise routes
// @title web-api
// @description The backend powering most things
// @contact.name API Support
// @contact.url https://github.com/ystv/web-api
// @contact.email computing@ystv.co.uk
func Init(version, commit string, db *sqlx.DB, cdn *s3.S3) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	debug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		panic(err)
	}
	// Enabling debugging
	e.Debug = debug

	creatorV1 := creatorPackage.NewRepos(db, cdn)
	miscV1 := miscPackage.NewRepos(db)
	publicV1 := publicPackage.NewRepos(db)
	peopleV1 := peoplePackage.NewRepo(db)
	clapperV1 := clapperPackage.NewRepos(db)

	// Authentication middleware
	middleware.Init(e)
	config := echoMw.JWTConfig{
		Claims:      &utils.JWTClaims{},
		TokenLookup: "cookie:token",
		SigningKey:  []byte(os.Getenv("signing_key")),
	}

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// List all possible routes
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
		// Service web endpoints
		encoder := internal.Group("/encoder")
		{
			encoder.POST("/upload_request", encoderV1.VideoNew)
		}
		stream := internal.Group("/stream")
		{
			stream.POST("/auth", streamV1.CheckAuth)
		}
		// Internal user endpoints
		if !debug {
			internal.Use(echoMw.JWTWithConfig(config))
		}
		{
			people := internal.Group("/people")
			{
				user := people.Group("/user")
				{
					user.GET("/full", peopleV1.UserByTokenFull)
					user.GET("/:id", peopleV1.UserByID)
					user.GET("/:id/full", peopleV1.UserByIDFull)
					user.GET("", peopleV1.UserByToken)
				}
				users := people.Group("/users")
				{
					users.GET("", peopleV1.ListAllPeople)
					users.GET("/:id", peopleV1.ListRoleMembers)
				}

			}
			creator := internal.Group("/creator")
			{
				videos := creator.Group("/videos")
				{
					videos.GET("", creatorV1.VideoList)
					videos.GET("/my", creatorV1.ListVideosByUser)
					videos.POST("", creatorV1.NewVideo)
					videoItem := videos.Group("/:id")
					{
						videoItem.GET("", creatorV1.GetVideo)
						videoItem.PUT("", notImplemented)
					}
				}
				series := creator.Group("/series")
				{
					series.GET("", creatorV1.ListSeries)
					series.GET("/:id", creatorV1.GetSeries)
				}
				playlists := creator.Group("/playlists")
				{
					playlists.GET("", creatorV1.ListPlaylist)
					playlists.POST("", creatorV1.NewPlaylist)
					playlist := playlists.Group("/:id")
					{
						playlist.GET("", creatorV1.GetPlaylist)
						playlist.PUT("", creatorV1.UpdatePlaylist)
					}
				}
				encodes := creator.Group("/encodes")
				{
					presets := encodes.Group("/presets")
					{
						presets.GET("", creatorV1.ListPreset)
						presets.POST("", creatorV1.NewPreset)
						presets.PUT("", creatorV1.UpdatePreset) // We take the ID in the json request
					}
					profiles := encodes.Group("/profiles")
					{
						profiles.GET("", creatorV1.ListEncodeProfile)
					}
				}
				creator.GET("/calendar/:year/:month", creatorV1.ListVideosByMonth)
				creator.GET("/stats", creatorV1.Stats)
			}
			clapper := internal.Group("/clapper")
			{
				calendar := clapper.Group("/calendar")
				{
					calendar.GET("/:year/:term", notImplemented)       // List all events of term
					calendar.GET("/:year/:month", clapperV1.ListMonth) // List all events of month
				}
				events := clapper.Group("/event")
				{
					events.POST("", clapperV1.NewEvent)   // Create a new event
					events.PUT("", clapperV1.UpdateEvent) // Update an event
					event := events.Group("/:eventid")
					{
						event.GET("", clapperV1.GetEvent) // Get event info, returns event info and signup sheets
						event.POST("/signup", clapperV1.NewSignup)
						signup := event.Group("/:signupid")
						{
							signup.PUT("", clapperV1.UpdateSignup)         // Create a new signup sheet
							signup.POST("/:positionid", clapperV1.NewCrew) // Add position to signup
							crew := event.Group("/:crewid")
							{
								crew.PUT("/reset", clapperV1.ResetCrew) // Set the role back to unassigned
								crew.PUT("", clapperV1.SetCrew)         // Update a crew role to the requesting user
								crew.DELETE("", clapperV1.DeleteCrew)   // Delete the crew role from signup
							}
						}
					}
				}
				positions := clapper.Group("/positions")
				{
					positions.GET("", clapperV1.ListPosition)   // List crew positions
					positions.POST("", clapperV1.NewPosition)   // Create a new crew position
					positions.PUT("", clapperV1.UpdatePosition) // Update a position
				}
			}
			misc := internal.Group("/misc")
			{
				quotes := misc.Group("/quotes")
				{
					quotes.GET("/:amount/:page", miscV1.ListQuotes)
					quotes.POST("", miscV1.NewQuote)
					quotes.PUT("", miscV1.UpdateQuote)
					quotes.DELETE("/:id", miscV1.DeleteQuote)
				}
				webcams := misc.Group("/webcams")
				{
					webcams.GET("/:id/*", miscV1.GetWebcam)
					webcams.GET("", miscV1.ListWebcams)
				}
			}
		}
		public := apiV1.Group("/public")
		{
			public.GET("/find/*", publicV1.Find)
			public.GET("/videos/:offset/:page", publicV1.ListVideos)
			public.GET("/video/:id", publicV1.Video)
			public.GET("/video/:id/breadcrumb", publicV1.VideoBreadcrumb)
			public.GET("/series/:id", publicV1.SeriesByID)
			public.GET("/series/:id/breadcrumb", publicV1.SeriesBreadcrumb)
			public.GET("/teams", publicV1.ListTeams)
			public.GET("/streams", publicPackage.StreamList)
			public.GET("/stream/:id", publicPackage.StreamFind)
			public.GET("/streams/home", publicPackage.StreamHome) // isLive null
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
