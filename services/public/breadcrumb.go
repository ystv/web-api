package public

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/ystv/web-api/utils"
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
	// BreadcrumbItem is either a video or a series
	BreadcrumbItem struct {
		Video  *VideoItem `json:"video,omitempty"`
		Series *Series    `json:"series,omitempty"`
	}
)

// SeriesBreadcrumb will return the breadcrumb from SeriesID to root
func SeriesBreadcrumb(seriesID int) ([]Breadcrumb, error) {
	s := []Breadcrumb{}
	// TODO Need a bool to indicate if series is in URL
	err := utils.DB.Select(&s,
		`SELECT parent.series_id as id, parent.url as url, COALESCE(parent.name, parent.url) as name
		FROM
			video.series node,
			video.series parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND node.series_id = $1
		ORDER BY parent.lft;`, seriesID)
	if err != nil {
		log.Printf("BreadcrumbSeries failed: %+v", err)
	}
	return s, err
}

// Find returns either a series or video for a given path
func Find(ctx context.Context, path string) (*BreadcrumbItem, error) {
	s, err := SeriesFromPath(path)
	if err != nil {
		// Might be a video, so we'll go one layer back and check for series
		if err == sql.ErrNoRows {
			split := strings.Split(path, "/")
			pathWithoutLast := strings.Join(split[:len(split)-1], "/")
			s, err := SeriesFromPath(pathWithoutLast)
			if err != nil {
				if err == sql.ErrNoRows {
					// No series, so there will be no videos
					return nil, errors.New("Breadcrumb: No series")
				}
				// log.Print(err)
				return nil, err
			}
			// Found series
			if len(s.ChildVideos) == 0 {
				// No videos on series
				return nil, errors.New("Series: No videos")
			}
			// We've got videos
			for _, v := range s.ChildVideos {
				// Check if video  name matches last name
				if v.URL == split[len(split)-1] {
					// Found video
					foundVideo, err := VideoFind(v.VideoID)
					if err != nil {
						return nil, err
					}
					return &BreadcrumbItem{foundVideo, nil}, nil
				}
			}
		} else {
			return nil, err
		}
	}
	// Found series
	return &BreadcrumbItem{nil, &s}, nil
}
