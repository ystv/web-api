package public

import (
	"context"

	"github.com/jmoiron/sqlx"
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
	}
	PlaylistRepo interface {
		GetPlaylist(ctx context.Context, playlistID int) (Playlist, error)
	}
	// BreadcrumbRepo represents all breadcrumb interactions
	BreadcrumbRepo interface {
		VideoBreadcrumb(ctx context.Context, videoID int) ([]Breadcrumb, error)
		SeriesBreadcrumb(ctx context.Context, seriesID int) ([]Breadcrumb, error)
		Find(ctx context.Context, path string) (*BreadcrumbItem, error)
	}
	// TeamRepo represents all team interactions
	TeamRepo interface {
		ListTeams(ctx context.Context) ([]Team, error)
		GetTeam(ctx context.Context, teamID int) (Team, error)
		GetTeamByYear(ctx context.Context, teamID, year int) (Team, error)
		ListTeamMembers(ctx context.Context, teamID int) ([]TeamMember, error)
		ListOfficers(ctx context.Context) ([]TeamMember, error)
	}
	// Store encapsulates our dependency
	Store struct {
		db *sqlx.DB
	}
)

// NewStore creates our data store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db}
}
