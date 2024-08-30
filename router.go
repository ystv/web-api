package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	clapperPackage "github.com/ystv/web-api/controllers/v1/clapper"
	creatorPackage "github.com/ystv/web-api/controllers/v1/creator"
	encoderPackage "github.com/ystv/web-api/controllers/v1/encoder"
	miscPackage "github.com/ystv/web-api/controllers/v1/misc"
	peoplePackage "github.com/ystv/web-api/controllers/v1/people"
	publicPackage "github.com/ystv/web-api/controllers/v1/public"
	streamV1 "github.com/ystv/web-api/controllers/v1/stream"
	"github.com/ystv/web-api/middleware"

	_ "github.com/ystv/web-api/swagger"

	echoSwagger "github.com/swaggo/echo-swagger"
	// Swag CLI generates documentation, you have to import it.
	// Run `go generate` if your IDE gives an import error here.
	//_ "github.com/ystv/web-api/swagger"
	"github.com/ystv/web-api/utils"
)

// TODO standardise on function names

// Router provides an HTTP server for web-api
type Router struct {
	version string
	commit  string
	router  *echo.Echo
	access  *utils.Accesser
	clapper *clapperPackage.Repos
	creator *creatorPackage.Repos
	encoder *encoderPackage.Repo
	misc    *miscPackage.Repos
	people  *peoplePackage.Repo
	public  *publicPackage.Repos
	stream  *streamV1.Repos
}

// NewRouter is the required dependencies
type NewRouter struct {
	Version    string
	Commit     string
	DomainName string
	Debug      bool
	Access     *utils.Accesser
	Clapper    *clapperPackage.Repos
	Creator    *creatorPackage.Repos
	Encoder    *encoderPackage.Repo
	Misc       *miscPackage.Repos
	People     *peoplePackage.Repo
	Public     *publicPackage.Repos
	Stream     *streamV1.Repos
}

// New creates a new router instance
func New(conf *NewRouter) *Router {
	r := &Router{
		version: conf.Version,
		commit:  conf.Commit,
		router:  echo.New(),
		access:  conf.Access,
		clapper: conf.Clapper,
		creator: conf.Creator,
		encoder: conf.Encoder,
		misc:    conf.Misc,
		people:  conf.People,
		public:  conf.Public,
		stream:  conf.Stream,
	}
	r.router.HideBanner = true

	// Enabling debugging
	r.router.Debug = conf.Debug

	// Authentication middleware
	middleware.New(r.router, conf.DomainName)

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
			encoder.POST("/upload_request", r.encoder.UploadRequest)
			encoder.POST("/transcode_finished/:taskid", r.encoder.TranscodeFinished)
		}
		stream := internal.Group("/stream")
		{
			stream.POST("/publish", r.stream.PublishHandler)
			stream.POST("/unpublish", r.stream.UnpublishHandler)
		}
		// Internal user endpoints
		if !r.router.Debug {
			internal.Use(r.access.AuthMiddleware)
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
					addUser := user.Group("/add")
					{
						if !r.router.Debug {
							addUser.Use(r.access.AddUserAuthMiddleware)
						}
						addUser.POST("", r.people.AddUser)
					}
					modifyUser := user.Group("/admin")
					{
						if !r.router.Debug {
							modifyUser.Use(r.access.ModifyUserAuthMiddleware)
						}
					}
				}
				users := people.Group("/users")
				{
					if !r.router.Debug {
						users.Use(r.access.ListUserAuthMiddleware)
					}
					users.GET("", r.people.ListAllPeople)
				}
				role := people.Group("/role")
				{
					role.GET("", r.people.ListAllRoles)
					role.GET("/:roleid/members", r.people.ListRoleMembersByID)
					role.GET("/:roleid/permissions", r.people.ListRolePermissionsByID)
				}
				permission := people.Group("/permission")
				{
					permission.GET("", r.people.ListAllPermissions)
					permission.GET("/:permissionid/members", r.people.ListPermissionMembersByID)
				}
			}
			creator := internal.Group("/creator")
			{
				videos := creator.Group("/video")
				{
					videos.GET("", r.creator.VideoList)
					videos.GET("/my", r.creator.ListVideosByUser)
					videos.POST("", r.creator.NewVideo)
					videos.PUT("/meta", r.creator.UpdateVideoMeta)
					videos.POST("/search", r.creator.SearchVideo)
					videoItem := videos.Group("/:id")
					{
						videoItem.GET("", r.creator.GetVideo)
						videoItem.DELETE("", r.creator.DeleteVideo)
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
				playlists := creator.Group("/playlist")
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
				playout := creator.Group("/playout")
				{
					playout.GET("/channels", r.creator.ListChannels)
					playout.POST("/channels", r.creator.NewChannel)
					playout.PUT("/channels", r.creator.UpdateChannel)
					playout.DELETE("/channels/:channelid", r.creator.DeleteChannel)
				}
				encode := creator.Group("/encode")
				{
					preset := encode.Group("/preset")
					{
						preset.GET("", r.creator.ListEncodePreset)
						preset.POST("", r.creator.NewEncodePreset)
						preset.PUT("", r.creator.UpdateEncodePreset) // We take the ID in the json request
						preset.DELETE("", r.creator.DeleteEncodePreset)
					}
					format := encode.Group("/format")
					{
						format.GET("", r.creator.ListEncodeFormat)
						format.PUT("", r.creator.UpdateEncodeFormat)
						format.POST("", r.creator.NewEncodeFormat)
						format.DELETE("/:formatid", r.creator.DeleteEncodeFormat)
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
					calendar.GET("/monthly/:year/:month", r.clapper.ListMonth) // List all events of the month
				}
				events := clapper.Group("/event")
				{
					events.POST("", r.clapper.NewEvent)   // Create a new event
					events.PUT("", r.clapper.UpdateEvent) // Update an event
					event := events.Group("/:eventid")
					{
						event.GET("", r.clapper.GetEvent) // Get event info, return event info and signup sheets
						event.DELETE("", r.clapper.DeleteEvent)
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
					positions.GET("", r.clapper.ListPosition)                  // List crew positions
					positions.POST("", r.clapper.NewPosition)                  // Create a new crew position
					positions.PUT("", r.clapper.UpdatePosition)                // Update a position
					positions.DELETE("/:positionid", r.clapper.DeletePosition) // Delete a position
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
				list := misc.Group("/list")
				{
					list.GET("s", r.misc.GetLists)
					list.GET("s/my", r.misc.GetListsByToken)
					list.GET("/:listid", r.misc.GetList)
					list.GET("/:listid/subscribers", r.misc.GetSubscribers)
					list.POST("/:listid/subscribe", r.misc.SubscribeByToken)
					list.POST("/:listid/subscribe/:userid", r.misc.SubscribeByID)
					list.DELETE("/:listid/unsubscribe", r.misc.UnsubscribeByToken)
					list.DELETE("/:listid/unsubscribe/:userid", r.misc.UnsubscribeByID)
				}
			}
		}
		apiV1.GET("/list-unsubscribe/:uuid", r.misc.UnsubscribeByUUID)

		public := apiV1.Group("/public")
		{
			public.POST("/search", r.public.Search)
			public.GET("/find/*", r.public.Find)
			video := public.Group("/video")
			{
				// /videos
				video.GET("s/:offset/:page", r.public.ListVideos)
				// /video
				video.GET("/:id", r.public.Video)
				video.GET("/:id/breadcrumb", r.public.VideoBreadcrumb)
			}
			series := public.Group("/series")
			{
				series.GET("/:id", r.public.SeriesByID)
				series.GET("/:id/breadcrumb", r.public.SeriesBreadcrumb)
				series.GET("/yearly/:year", r.public.SeriesByYear)
			}
			playlist := public.Group("/playlist")
			{
				playlist.GET("/random", r.public.GetPlaylistRandom)
				playlist.GET("/popular/all", r.public.GetPlaylistPopularByAllTime)
				playlist.GET("/popular/year", r.public.GetPlaylistPopularByPastYear)
				playlist.GET("/popular/month", r.public.GetPlaylistPopularByPastMonth)
				playlist.GET("/:playlistid", r.public.GetPlaylist)
			}
			teams := public.Group("/teams")
			{
				teams.GET("", r.public.ListTeams)
				teams.GET("/officers", r.public.ListOfficers)
				teamsEmail := teams.Group("/email")
				{
					teamsEmail.GET("/:emailAlias/:startYear/:endYear", r.public.GetTeamByStartEndYearByEmail)
					teamsEmail.GET("/:emailAlias/:year", r.public.GetTeamByYearByEmail)
					teamsEmail.GET("/:emailAlias", r.public.GetTeamByEmail)
				}
				teamsId := teams.Group("/teamid")
				{
					teamsId.GET("/:teamid/:startYear/:endYear", r.public.GetTeamByStartEndYearById)
					teamsId.GET("/:teamid/:year", r.public.GetTeamByYearById)
					teamsId.GET("/:teamid", r.public.GetTeamById)
				}
			}
			streamChannel := public.Group("/playout/channel")
			{
				streamChannel.GET("", r.public.ListChannels)
				streamChannel.GET("/:channelShortName", r.public.GetChannel)
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
*/
