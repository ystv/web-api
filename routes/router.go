package routes

import (
	"fmt"
	"net/http"
	"time"

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

// Router provides a HTTP server for web-api
type Router struct {
	version   string
	commit    string
	router    *echo.Echo
	jwtConfig *echoMw.JWTConfig
	clapper   *clapperPackage.Repos
	creator   *creatorPackage.Repos
	misc      *miscPackage.Repos
	people    *peoplePackage.Repo
	public    *publicPackage.Repos
}

// NewRouter is the required dependencies
type NewRouter struct {
	Version       string
	Commit        string
	JWTSigningKey string
	Debug         bool
	Clapper       *clapperPackage.Repos
	Creator       *creatorPackage.Repos
	Misc          *miscPackage.Repos
	People        *peoplePackage.Repo
	Public        *publicPackage.Repos
}

// New creates a new router instance
func New(conf *NewRouter) *Router {
	r := &Router{
		version: conf.Version,
		commit:  conf.Commit,
		router:  echo.New(),
		jwtConfig: &echoMw.JWTConfig{
			Claims:      &utils.JWTClaims{},
			TokenLookup: "cookie:token",
			SigningKey:  []byte(conf.JWTSigningKey),
		},
		clapper: conf.Clapper,
		creator: conf.Creator,
		misc:    conf.Misc,
		people:  conf.People,
		public:  conf.Public,
	}
	r.router.HideBanner = true

	// Enabling debugging
	r.router.Debug = conf.Debug

	// Authentication middleware
	middleware.Init(r.router)

	r.loadRoutes()

	return r
}

// Start the HTTP Server
func (r *Router) Start() {
	r.router.Logger.Fatal(r.router.Start(":8081"))
}

// loadRoutes initialise routes
// @title web-api
// @description The backend powering most things
// @contact.name API Support
// @contact.url https://github.com/ystv/web-api
// @contact.email computing@ystv.co.uk
func (r *Router) loadRoutes() {
	// swagger
	r.router.GET("/swagger/*", echoSwagger.WrapHandler)

	// List all possible routes
	r.router.GET("/routes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, r.router.Routes())
	})
	// ping
	r.router.GET("/ping", func(c echo.Context) error {
		resp := map[string]time.Time{"pong": time.Now()}
		return c.JSON(http.StatusOK, resp)
	})

	apiV1 := r.router.Group("v1")
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
		if !r.router.Debug {
			internal.Use(echoMw.JWTWithConfig(*r.jwtConfig))
		}
		{
			people := internal.Group("/people")
			{
				user := people.Group("/user")
				{
					user.GET("/full", r.people.UserByTokenFull)
					user.GET("/:id", r.people.UserByID)
					user.GET("/:id/full", r.people.UserByIDFull)
					user.GET("", r.people.UserByToken)
				}
				users := people.Group("/users")
				{
					users.GET("", r.people.ListAllPeople)
					users.GET("/:id", r.people.ListRoleMembers)
				}

			}
			creator := internal.Group("/creator")
			{
				videos := creator.Group("/videos")
				{
					videos.GET("", r.creator.VideoList)
					videos.GET("/my", r.creator.ListVideosByUser)
					videos.POST("", r.creator.NewVideo)
					videoItem := videos.Group("/:id")
					{
						videoItem.GET("", r.creator.GetVideo)
						videoItem.PUT("", notImplemented)
						// videoItem.DELETE("", r.creator.DeleteVideo)
					}
				}
				series := creator.Group("/series")
				{
					series.GET("", r.creator.ListSeries)
					seriesItem := series.Group("/:seriesid")
					{
						seriesItem.GET("", r.creator.GetSeries)
						// seriesItem.PUT("", r.creator.UpdateSeries)
						// seriesItem.DELETE("", r.creator.DeleteSeries)
					}
				}
				playlists := creator.Group("/playlists")
				{
					playlists.GET("", r.creator.ListPlaylist)
					playlists.POST("", r.creator.NewPlaylist)
					playlist := playlists.Group("/:id")
					{
						playlist.GET("", r.creator.GetPlaylist)
						playlist.PUT("", r.creator.UpdatePlaylist)
						// playlist.DELETE("", r.creator.DeletePlaylist)
					}
				}
				encodes := creator.Group("/encodes")
				{
					presets := encodes.Group("/presets")
					{
						presets.GET("", r.creator.ListPreset)
						presets.POST("", r.creator.NewPreset)
						presets.PUT("", r.creator.UpdatePreset) // We take the ID in the json request
					}
					profiles := encodes.Group("/profiles")
					{
						profiles.GET("", r.creator.ListEncodeProfile)
					}
				}
				creator.GET("/calendar/:year/:month", r.creator.ListVideosByMonth)
				creator.GET("/stats", r.creator.Stats)
			}
			clapper := internal.Group("/clapper")
			{
				calendar := clapper.Group("/calendar")
				{
					calendar.GET("/termly/:year/:term", notImplemented)        // List all events of term
					calendar.GET("/monthly/:year/:month", r.clapper.ListMonth) // List all events of month
				}
				events := clapper.Group("/event")
				{
					events.POST("", r.clapper.NewEvent)   // Create a new event
					events.PUT("", r.clapper.UpdateEvent) // Update an event
					event := events.Group("/:eventid")
					{
						event.GET("", r.clapper.GetEvent) // Get event info, returns event info and signup sheets
						event.POST("/signup", r.clapper.NewSignup)
						signup := event.Group("/:signupid")
						{
							signup.PUT("", r.clapper.UpdateSignup)         // Create a new signup sheet
							signup.POST("/:positionid", r.clapper.NewCrew) // Add position to signup
							crew := event.Group("/:crewid")
							{
								crew.PUT("/reset", r.clapper.ResetCrew) // Set the role back to unassigned
								crew.PUT("", r.clapper.SetCrew)         // Update a crew role to the requesting user
								crew.DELETE("", r.clapper.DeleteCrew)   // Delete the crew role from signup
							}
						}
					}
				}
				positions := clapper.Group("/positions")
				{
					positions.GET("", r.clapper.ListPosition)   // List crew positions
					positions.POST("", r.clapper.NewPosition)   // Create a new crew position
					positions.PUT("", r.clapper.UpdatePosition) // Update a position
				}
			}
			misc := internal.Group("/misc")
			{
				quotes := misc.Group("/quotes")
				{
					quotes.GET("/:amount/:page", r.misc.ListQuotes)
					quotes.POST("", r.misc.NewQuote)
					quotes.PUT("", r.misc.UpdateQuote)
					quotes.DELETE("/:id", r.misc.DeleteQuote)
				}
				webcams := misc.Group("/webcams")
				{
					webcams.GET("/:id/*", r.misc.GetWebcam)
					webcams.GET("", r.misc.ListWebcams)
				}
			}
		}
		public := apiV1.Group("/public")
		{
			public.GET("/find/*", r.public.Find)
			video := public.Group("/video")
			{
				video.GET("s/:offset/:page", r.public.ListVideos)
				video.GET("/:id", r.public.Video)
				video.GET("/:id/breadcrumb", r.public.VideoBreadcrumb)
			}
			series := public.Group("/series")
			{
				series.GET("/:id", r.public.SeriesByID)
				series.GET("/:id/breadcrumb", r.public.SeriesBreadcrumb)
			}
			public.GET("/playlist/:playlistid", r.public.GetPlaylist)
			teams := public.Group("/teams")
			{
				teams.GET("", r.public.ListTeams)
				teams.GET("/officers", r.public.ListOfficers)
				teams.GET("/:teamid", r.public.GetTeam)
				teams.GET("/:teamid/:year", r.public.GetTeamByYear)
			}
			streams := public.Group("/streams")
			{
				streams.GET("", publicPackage.StreamList)
				streams.GET("/:id", publicPackage.StreamFind)
				streams.GET("/home", publicPackage.StreamHome) // isLive null
			}

		}

	}
	r.router.GET("/", func(c echo.Context) error {
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
`, r.version, r.commit)
		return c.String(http.StatusOK, text)
	})
}

func notImplemented(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

/*
- by year
- popular
- recent
- genre per
- featured playlist
-
*/
