package public

import (
	"context"
	"time"

	"github.com/jackc/pgx"
)

type (
	// VideoRepo represents all video interactions
	VideoRepo interface {
		ListVideo(ctx context.Context, offset int, page int) (*[]VideoMeta, error)
		GetVideo(ctx context.Context, videoID int) (*VideoItem, error)
		VideoOfSeries(ctx context.Context, seriesID int) ([]VideoMeta, error)
	}
	// SeriesRepo represents all series interactions
	SeriesRepo interface {
		GetSeries(ctx context.Context, seriesID int) (Series, error)
		GetSeriesMeta(ctx context.Context, seriesID int) (Series, error)
		GetSeriesImmediateChildrenSeries(ctx context.Context, seriesID int) ([]SeriesMeta, error)
		GetSeriesFromPath(ctx context.Context, path string) (Series, error)
		Search(ctx context.Context, query string) (Series, error)
	}
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
		VideoBreadcrumb(ctx context.Context, videoID int) ([]Breadcrumb, error)
		SeriesBreadcrumb(ctx context.Context, seriesID int) ([]Breadcrumb, error)
		Find(ctx context.Context, path string) (BreadcrumbItem, error)
	}
	// TeamRepo represents all team interactions
	TeamRepo interface {
		ListTeams(ctx context.Context) ([]Team, error)
		GetTeam(ctx context.Context, teamID int) (Team, error)
		GetTeamByYear(ctx context.Context, teamID, year int) (Team, error)
		ListTeamMembers(ctx context.Context, teamID int) ([]TeamMember, error)
		ListOfficers(ctx context.Context) ([]TeamMember, error)
	}
	// StreamRepo represents all stream / playout interactions
	StreamRepo interface {
		ListChannels(ctx context.Context) ([]Channel, error)
		GetChannel(ctx context.Context, urlName string) (Channel, error)
	}
	// Store encapsulates our dependency
	Store struct {
		db *pgx.Conn
	}
)

// NewStore creates our data store
func NewStore(db *pgx.Conn) *Store {
	return &Store{db}
}
