package creator

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/creator/types/breadcrumb"
	"github.com/ystv/web-api/services/creator/types/encode"
	"github.com/ystv/web-api/services/creator/types/playlist"
	"github.com/ystv/web-api/services/creator/types/playout"
	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/types/stats"
	"github.com/ystv/web-api/services/creator/types/video"
)

type (
	Repo interface {
		StatRepo
	}

	// Config configures where creator will use as its bucket sources
	Config struct {
		IngestBucket string
		ServeBucket  string
		Endpoint     string
	}
	// VideoRepo defines all creator video interactions
	VideoRepo interface {
		// GetItem gets the individual video item
		GetItem(ctx context.Context, id int) (video.ItemDB, error)
		ListMeta(ctx context.Context) ([]video.MetaDB, error)
		ListMetaByUser(ctx context.Context, userID int) ([]video.MetaDB, error)
		ListByCalendarMonth(ctx context.Context, year, month int) ([]video.MetaCal, error)
		OfSeries(ctx context.Context, seriesID int) ([]video.MetaDB, error)
		Search(ctx context.Context, query string) ([]video.MetaDB, error)
		// NewItem inserts a new video
		NewItem(ctx context.Context, v video.New) (int, error)
		// UpdateMeta updates the video metadata
		UpdateMeta(ctx context.Context, meta video.Meta) error
		// DeleteItem removes a video
		DeleteItem(ctx context.Context, videoID, userID int) error
		DeleteItemPermanently(ctx context.Context, videoID int) error
		// DeleteFile(ctx context.Context, fileID, userID int) error
	}
	// SeriesRepo defines all creator series interactions
	SeriesRepo interface {
		GetSeries(ctx context.Context, seriesID int) (series.SeriesDB, error)
		GetMeta(ctx context.Context, seriesID int) (series.Meta, error)
		ImmediateChildrenSeries(ctx context.Context, seriesID int) ([]series.Meta, error)
		List(ctx context.Context) ([]series.Meta, error)
		FromPath(ctx context.Context, path string) (series.SeriesDB, error)
	}
	// ChannelRepo defines all channel interactions
	ChannelRepo interface {
		ListChannels(ctx context.Context) ([]playout.Channel, error)
		GetChannel(ctx context.Context, urlName string) (playout.Channel, error)
		NewChannel(ctx context.Context, ch playout.Channel) error
		UpdateChannel(ctx context.Context, ch playout.Channel) error
		DeleteChannel(ctx context.Context, urlName string) error
	}
	// PlaylistRepo defines all playlist interactions
	PlaylistRepo interface {
		ListPlaylists(ctx context.Context) ([]playlist.PlaylistDB, error)
		GetPlaylist(ctx context.Context, playlistID int) (playlist.PlaylistDB, error)
		NewPlaylist(ctx context.Context, p playlist.New) (int, error)
		UpdatePlaylist(ctx context.Context, p playlist.Meta, videoIDs []int) error
		DeletePlaylist(ctx context.Context, playlistID int) error
		AddVideo(ctx context.Context, playlistID, videoID int) error
		DeleteVideo(ctx context.Context, playlistID, videoID int) error
		AddVideos(ctx context.Context, playlistID int, videoIDs []int) error
	}
	// BreadcrumbRepo defines all creator breadcrumb interactions
	BreadcrumbRepo interface {
		Series(ctx context.Context, seriesID int) ([]breadcrumb.Breadcrumb, error)
		Video(ctx context.Context, videoID int) ([]breadcrumb.Breadcrumb, error)
		Find(ctx context.Context, path string) (breadcrumb.Item, error)
	}
	// EncodeRepo defines all encode interactions
	EncodeRepo interface {
		ListFormat(ctx context.Context) ([]encode.Format, error)
		NewFormat(ctx context.Context, format encode.Format) (int, error)
		UpdateFormat(ctx context.Context, format encode.Format) error
		DeleteFormat(ctx context.Context, formatID int) error
		GetPreset(ctx context.Context, presetID int) (encode.Preset, error)
		ListPreset(ctx context.Context) ([]encode.Preset, error)
		NewPreset(ctx context.Context, p encode.Preset) (int, error)
		UpdatePreset(ctx context.Context, p encode.Preset) error
		DeletePreset(ctx context.Context, presetID int) error
	}
	// StatRepo defines all statistical interactions
	StatRepo interface {
		GlobalVideoStats(ctx context.Context) (stats.VideoGlobalStats, error)
	}
)

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) Repo {
	return &Store{
		db: db,
	}
}
