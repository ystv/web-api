package series

import (
	"github.com/ystv/web-api/services/creator/types/video"
	"gopkg.in/guregu/null.v4"
)

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see it's immediate children.
	Series struct {
		*Meta
		ImmediateChildSeries *[]Meta       `json:"childSeries"`
		ChildVideos          *[]video.Meta `json:"videos"`
	}
	// Meta is used as a children object for a series
	Meta struct {
		SeriesID    int         `db:"series_id" json:"id"`
		URL         string      `db:"url" json:"url"`
		SeriesName  null.String `db:"name" json:"name"`
		Description null.String `db:"description" json:"description"`
		Thumbnail   null.String `db:"thumbnail" json:"thumbnail"`
		Depth       int         `db:"depth" json:"depth"`
	}
)
