package breadcrumb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

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

// Item is either a video or a series
type Item struct {
	Video  *video.Item
	Series series.Series
}

// Find will returns either a series or a video for a given path
func Find(path string) (Item, error) {
	blank := Item{}

	Series, err := series.FromPath(path)
	if err != nil {
		// Might be a video, so we'll go back a crumb and check for a series
		if err == sql.ErrNoRows {
			split := strings.Split(path, "/")
			PathWithoutLast := strings.Join(split[:len(split)-1], "/")
			Series, err := series.FromPath(PathWithoutLast)
			if err != nil {
				if err == sql.ErrNoRows {
					// No series, so there will be no videos
					return blank, err
				}
				log.Printf("Find failed from 2nd last: %+v", err)
				return blank, err
			}
			// Found series
			if len(Series.ChildVideos) == 0 {
				// No videos on series
				return blank, errors.New("Series: No videos")
			}
			// We've got videos
			for _, v := range Series.ChildVideos {
				// Check if video name matches last path
				if v.URL == split[len(split)-1] {
					// Found video
					foundVideo, _ := video.FindVideoItem(context.Background(), v.ID)
					return Item{Video: foundVideo}, nil

				}
			}
		} else {
			log.Printf("Find failed from path: %+v", err)
			return blank, err
		}
	}
	// Found series
	return Item{Series: Series}, nil
}
