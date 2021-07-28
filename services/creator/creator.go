package creator

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/creator/types/breadcrumb"
	"github.com/ystv/web-api/services/creator/types/encode"
	"github.com/ystv/web-api/services/creator/types/playlist"
	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/types/stats"
	"github.com/ystv/web-api/services/creator/types/video"
)

type (
	// Config configures where creator will use as its bucket sources
	Config struct {
		IngestBucket string
		ServeBucket  string
	}
	// VideoRepo defines all creator video interactions
	VideoRepo interface {
		// Get methods
		GetItem(ctx context.Context, id int) (*video.Item, error)
		ListMeta(ctx context.Context) (*[]video.Meta, error)
		ListMetaByUser(ctx context.Context, userID int) (*[]video.Meta, error)
		ListByCalendarMonth(ctx context.Context, year, month int) (*[]video.MetaCal, error)
		OfSeries(ctx context.Context, seriesID int) (*[]video.Meta, error)

		// New method
		NewItem(ctx context.Context, v *video.New) (int, error)

		// Update method
		UpdateMeta(ctx context.Context, meta video.Meta) error

		// Delete method
		DeleteItem(ctx context.Context, videoID, userID int) error
		// DeleteFile(ctx context.Context, fileID, userID int) error
	}
	// SeriesRepo defines all creator series interactions
	SeriesRepo interface {
		Get(ctx context.Context, seriesID int) (*series.Series, error)
		GetMeta(ctx context.Context, seriesID int) (*series.Meta, error)
		ImmediateChildrenSeries(ctx context.Context, seriesID int) (*[]series.Meta, error)
		List(ctx context.Context) (*[]series.Meta, error)
		FromPath(ctx context.Context, path string) (*series.Series, error)
	}
	// BreadcrumbRepo defines all creator breadcrumb interactions
	BreadcrumbRepo interface {
		Series(ctx context.Context, seriesID int) (*[]breadcrumb.Breadcrumb, error)
		Video(ctx context.Context, videoID int) (*[]breadcrumb.Breadcrumb, error)
		Find(ctx context.Context, path string) (*breadcrumb.Item, error)
	}
	// PlaylistRepo defines all playlist interactions
	PlaylistRepo interface {
		All(ctx context.Context) ([]playlist.Playlist, error)
		Get(ctx context.Context, playlistID int) (playlist.Playlist, error)
		New(ctx context.Context, p playlist.New) (int, error)
		Update(ctx context.Context, p playlist.Meta, videoIDs []int) error
		AddVideo(ctx context.Context, playlistID, videoID int) error
		DeleteVideo(ctx context.Context, playlistID, videoID int) error
		AddVideos(ctx context.Context, playlistID int, videoIDs []int) error
	}
	// EncodeRepo defines all encode interactions
	EncodeRepo interface {
		GetPreset(ctx context.Context, presetID int) (*encode.Preset, error)
		ListPreset(ctx context.Context) ([]encode.Preset, error)
		NewPreset(ctx context.Context, p *encode.Preset) (int, error)
		UpdatePreset(ctx context.Context, p *encode.Preset) error
		ListFormat(ctx context.Context) ([]encode.Format, error)
	}
	// StatRepo defines all statistical interactions
	StatRepo interface {
		GlobalVideo(ctx context.Context) (*stats.VideoGlobalStats, error)
	}
)

// Here for validation to ensure we are meeting the interface
var _ StatRepo = &Store{}

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}
