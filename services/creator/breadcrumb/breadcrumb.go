package breadcrumb

import (
	"log"

	"github.com/ystv/web-api/services/creator/series"
	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/utils"
)

// Breadcrumb generic to be used for both series and video as a breadcrumb
type Breadcrumb struct {
	ID       int    `db:"id" json:"id"`
	URL      string `db:"url" json:"url"`
	UseInURL bool   `db:"use" json:"useInURL"`
	Name     string `db:"name" json:"name"`
	SeriesID int    `db:"series_id" json:"-"` // Here since needed
}

// Series will return the breadcrumb from SeriesID to root
func Series(SeriesID int) ([]Breadcrumb, error) {
	s := []Breadcrumb{}
	// TODO Need a bool to indicate if series is in URL
	err := utils.DB.Select(&s,
		`SELECT parent.series_id as id, parent.url as url, COALESCE(parent.name, parent.url) as name
		FROM
			video.series node,
			video.series parent
		WHERE
			node.series_left BETWEEN parent.series_left AND parent.series_right
			AND node.series_id = $1
		ORDER BY parent.series_left;`, SeriesID)
	if err != nil {
		log.Printf("BreadcrumbSeries failed: %+v", err)
	}
	return s, err
}

// Video returns the absolute path from a VideoID
func Video(VideoID int) ([]Breadcrumb, error) {
	var vB Breadcrumb // Video breadcrumb
	err := utils.DB.Get(&vB,
		`SELECT video_id as id, series_id, COALESCE(name, url) as name, url
		FROM video.items
		WHERE video_id = $1`, VideoID)
	if err != nil {
		log.Printf("VideoBreadcrumb failed: %+v", err)
		return nil, err
	}
	sB, err := Series(vB.SeriesID)
	if err != nil {
		return nil, err
	}
	sB = append(sB, vB)

	return sB, err
}

type Path struct {
	video.Item
	series.Series
}

func Find(path string)
