package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

var (
	ErrVideoNotFound  = errors.New("video not found")
	ErrSeriesNotFound = errors.New("series not found")
)

var _ BreadcrumbRepo = &Store{}

// VideoBreadcrumb returns the absolute path from a VideoID
func (s *Store) VideoBreadcrumb(ctx context.Context, videoID int) ([]Breadcrumb, error) {
	var vB Breadcrumb // Video breadcrumb
	err := s.db.GetContext(ctx, &vB,
		`SELECT video_id as id, series_id, COALESCE(name, url) as name, url
		FROM video.items
		WHERE video_id = $1`, videoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVideoNotFound
		}
		return nil, fmt.Errorf("failed to get video breadcrumb: %w", err)
	}
	sB, err := s.SeriesBreadcrumb(ctx, vB.SeriesID)
	if err != nil {
		// Interesting edge-case
		if !errors.Is(err, ErrSeriesNotFound) {
			return nil, fmt.Errorf("failed to get series breadcrumb: %w", err)
		}
	}
	sB = append(sB, vB)

	return sB, nil
}

// SeriesBreadcrumb will return the breadcrumb from SeriesID to root
func (s *Store) SeriesBreadcrumb(ctx context.Context, seriesID int) ([]Breadcrumb, error) {
	var b []Breadcrumb
	// TODO Need a bool to indicate if series is in URL
	err := s.db.SelectContext(ctx, &b,
		`SELECT parent.series_id as id, parent.url as url, COALESCE(parent.name, parent.url) as name
		FROM
			video.series node,
			video.series parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND parent.status = 'public'
			AND node.series_id = $1
		ORDER BY parent.lft;`, seriesID)
	if err != nil {
		return []Breadcrumb{}, ErrSeriesNotFound
	}
	// For some reason it's not returning NoRows
	if len(b) == 0 {
		return []Breadcrumb{}, ErrSeriesNotFound
	}
	return b, err
}

// Find returns either a series or video for a given path
// TODO be consistent with creator's find in terms of variables
func (s *Store) Find(ctx context.Context, path string) (BreadcrumbItem, error) {
	// Check to see if it's just a video ID
	videoID, err := strconv.Atoi(path)
	if err == nil {
		// It's a raw video ID
		foundVideo, err := s.GetVideo(ctx, videoID)
		if err == nil {
			return BreadcrumbItem{Video: foundVideo}, nil
		} else if errors.Is(err, sql.ErrNoRows) {
			return BreadcrumbItem{}, ErrVideoNotFound
		} else {
			return BreadcrumbItem{}, fmt.Errorf("failed to get video: %w", err)
		}
	}
	series, err := s.GetSeriesFromPath(ctx, path)
	if err != nil {
		// Might be a video, so we'll go one layer back and check for series
		if errors.Is(err, sql.ErrNoRows) {
			split := strings.Split(path, "/")
			pathWithoutLast := strings.Join(split[:len(split)-1], "/")
			series, err = s.GetSeriesFromPath(ctx, pathWithoutLast)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// No series, so there will be no videos
					return BreadcrumbItem{}, ErrSeriesNotFound
				}
				return BreadcrumbItem{}, err
			}
			// Found series, let's check for the video
			if len(series.ChildVideos) == 0 {
				// No videos on series
				return BreadcrumbItem{}, ErrVideoNotFound
			}
			// We've got videos
			for _, v := range series.ChildVideos {
				// Check if video name matches last name
				if v.URL == split[len(split)-1] {
					// Found video
					foundVideo, err := s.GetVideo(ctx, v.VideoID)
					if err != nil {
						return BreadcrumbItem{}, fmt.Errorf("failed to get video: %w", err)
					}
					return BreadcrumbItem{foundVideo, nil}, nil
				}
			}
		} else {
			return BreadcrumbItem{}, fmt.Errorf("failed to get series from path: %w", err)
		}
	}
	// Found series
	return BreadcrumbItem{nil, &series}, nil
}
