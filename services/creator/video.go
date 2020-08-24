package creator

import (
	"context"

	"github.com/ystv/web-api/services/creator/types/breadcrumb"
	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/types/video"
)

type (
	// VideoRepo defines all creator video interactions
	VideoRepo interface {
		GetItem(ctx context.Context, id int) (*video.Item, error)
		ListMeta(ctx context.Context) (*[]video.Meta, error)
		ListMetaByUser(ctx context.Context, userID int) (*[]video.Meta, error)
		ListByCalendarMonth(ctx context.Context, year, month int) (*[]video.MetaCal, error)
		OfSeries(ctx context.Context, seriesID int) (*[]video.Meta, error)
		NewItem(ctx context.Context, v *video.NewVideo) error
	}
	// SeriesRepo defines all creator series interactions
	SeriesRepo interface {
		View(ctx context.Context, seriesID int) (*series.Series, error)
		Info(ctx context.Context, seriesID int) (*series.Meta, error)
		ImmediateChildrenSeries(ctx context.Context, seriesID int) (*[]series.Meta, error)
		All(ctx context.Context) (*[]series.Meta, error)
		FromPath(ctx context.Context, path string) (*series.Series, error)
	}
	// BreadcrumbRepo defines all creator breadcrumb interactions
	BreadcrumbRepo interface {
		Series(ctx context.Context, seriesID int) (*[]breadcrumb.Breadcrumb, error)
		Video(ctx context.Context, videoID int) (*[]breadcrumb.Breadcrumb, error)
		Find(ctx context.Context, path string) (*breadcrumb.Item, error)
	}
)
