package public

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	Repos interface {
		VideoRepo
		SeriesRepo
		PlaylistRepo
		BreadcrumbRepo
		TeamRepo
		StreamRepo
		CustomSettingsRepo
	}

	// VideoRepo represents all video interactions
	VideoRepo interface {
		ListVideos(ctx context.Context, offset int, page int) (*[]VideoMeta, error)
		GetVideo(ctx context.Context, videoID int) (*VideoItem, error)
		VideoOfSeries(ctx context.Context, seriesID int) ([]VideoMeta, error)
	}
	// SeriesRepo represents all series interactions
	SeriesRepo interface {
		GetSeries(ctx context.Context, seriesID int) (Series, error)
		GetSeriesMeta(ctx context.Context, seriesID int) (SeriesMeta, error)
		GetSeriesFullMeta(ctx context.Context, seriesID int) (Series, error)
		GetSeriesImmediateChildrenSeries(ctx context.Context, seriesID int) ([]SeriesMeta, error)
		GetSeriesFromPath(ctx context.Context, path string) (Series, error)
		GetSeriesByYear(ctx context.Context, year int) (Series, error)
		Search(ctx context.Context, query string) (Series, error)
	}
	// PlaylistRepo represents all playlist interactions
	PlaylistRepo interface {
		GetPlaylist(ctx context.Context, playlistID int) (Playlist, error)
		GetPlaylistPopular(ctx context.Context, fromPeriod time.Time) (Playlist, error)
		GetPlaylistPopularByAllTime(ctx context.Context) (Playlist, error)
		GetPlaylistPopularByPastYear(ctx context.Context) (Playlist, error)
		GetPlaylistPopularByPastMonth(ctx context.Context) (Playlist, error)
		GetPlaylistRandom(ctx context.Context) (Playlist, error)
	}
	// BreadcrumbRepo represents all breadcrumb interactions
	BreadcrumbRepo interface {
		GetVideoBreadcrumb(ctx context.Context, videoID int) ([]Breadcrumb, error)
		GetSeriesBreadcrumb(ctx context.Context, seriesID int) ([]Breadcrumb, error)
		Find(ctx context.Context, path string) (BreadcrumbItem, error)
	}
	// TeamRepo represents all team interactions
	TeamRepo interface {
		ListTeams(ctx context.Context) ([]Team, error)
		GetTeamByEmail(ctx context.Context, emailAlias string) (Team, error)
		GetTeamByID(ctx context.Context, teamID int) (Team, error)
		GetTeamByYearByEmail(ctx context.Context, emailAlias string, year int) (Team, error)
		GetTeamByYearByID(ctx context.Context, teamID, year int) (Team, error)
		GetTeamByStartEndYearByEmail(ctx context.Context, emailAlias string, startYear, endYear int) (Team, error)
		GetTeamByStartEndYearByID(ctx context.Context, teamID, startYear, endYear int) (Team, error)
		ListTeamMembers(ctx context.Context, teamID int) ([]TeamMemberDB, error)
		ListOfficers(ctx context.Context) ([]TeamMemberDB, error)
	}
	// StreamRepo represents all stream / playout interactions
	StreamRepo interface {
		ListChannels(ctx context.Context) ([]Channel, error)
		GetChannel(ctx context.Context, urlName string) (Channel, error)
	}
	CustomSettingsRepo interface {
		GetCustomSettingPublic(ctx context.Context, settingID string) (CustomSetting, error)
	}
	// Store encapsulates our dependency
	Store struct {
		db          *sqlx.DB
		cdnEndpoint string
	}
)

// NewStore creates our data store
func NewStore(db *sqlx.DB, cdnEndpoint string) Repos {
	return &Store{
		db:          db,
		cdnEndpoint: cdnEndpoint,
	}
}
