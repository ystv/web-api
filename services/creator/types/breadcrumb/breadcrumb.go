package breadcrumb

import (
	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/types/video"
)

type (
	// Breadcrumb generic to be used for both series and video as a breadcrumb
	Breadcrumb struct {
		ID       int    `db:"id" json:"id"`
		URL      string `db:"url" json:"url"`
		UseInURL bool   `db:"use" json:"useInURL"`
		Name     string `db:"name" json:"name"`
		SeriesID int    `db:"series_id" json:"-"` // Here since needed
	}
	// Item is either a video or a series
	Item struct {
		Video  *video.Item    `json:"video"`
		Series *series.Series `json:"series"`
	}
)
