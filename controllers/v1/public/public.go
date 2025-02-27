package public

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/public"
)

// Repos encapsulates the dependency
type (
	Repos interface {
		BreadcrumbRepo
		PlaylistRepo
		SeriesRepo
		StreamRepo
		TeamRepo
		VideoRepo
	}

	BreadcrumbRepo interface {
		Find(c echo.Context) error
		VideoBreadcrumb(c echo.Context) error
		GetSeriesBreadcrumb(c echo.Context) error
	}

	PlaylistRepo interface {
		GetPlaylist(c echo.Context) error
		GetPlaylistPopularByAllTime(c echo.Context) error
		GetPlaylistPopularByPastYear(c echo.Context) error
		GetPlaylistPopularByPastMonth(c echo.Context) error
		GetPlaylistRandom(c echo.Context) error
	}

	SeriesRepo interface {
		GetSeriesByID(c echo.Context) error
		GetSeriesByYear(c echo.Context) error
		Search(c echo.Context) error
	}

	StreamRepo interface {
		ListChannels(c echo.Context) error
		GetChannel(c echo.Context) error
	}

	TeamRepo interface {
		ListTeams(c echo.Context) error
		GetTeamByEmail(c echo.Context) error
		GetTeamByID(c echo.Context) error
		GetTeamByYearByEmail(c echo.Context) error
		GetTeamByYearByID(c echo.Context) error
		GetTeamByStartEndYearByEmail(c echo.Context) error
		GetTeamByStartEndYearByID(c echo.Context) error
		ListOfficers(c echo.Context) error
	}

	VideoRepo interface {
		GetVideo(c echo.Context) error
		ListVideos(c echo.Context) error
	}

	Store struct {
		public public.Repos
	}
)

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) Repos {
	return &Store{public.NewStore(db)}
}
