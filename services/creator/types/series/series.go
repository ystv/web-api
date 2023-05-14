package series

import (
	"errors"

	"github.com/ystv/web-api/services/creator/types/video"
)

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see its immediate children.
	Series struct {
		Meta
		ImmediateChildSeries []Meta       `json:"childSeries"`
		ChildVideos          []video.Meta `json:"videos"`
	}
	// Meta is used as a children object for a series
	Meta struct {
		SeriesID    int    `db:"series_id" json:"id"`
		URL         string `db:"url" json:"url"`
		SeriesName  string `db:"name" json:"name"`
		Description string `db:"description" json:"description"`
		Thumbnail   string `db:"thumbnail" json:"thumbnail"`
		Depth       int    `db:"depth" json:"depth"`
	}
)

var (
	ErrNotFound               = errors.New("series not found")
	ErrMetaNotFound           = errors.New("series meta not found")
	ErrChildrenSeriesNotFound = errors.New("series children series not found")
	ErrChildrenVideosNotFound = errors.New("series videos not found")
)
