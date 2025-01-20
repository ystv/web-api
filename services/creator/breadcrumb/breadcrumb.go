package breadcrumb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/breadcrumb"
	"github.com/ystv/web-api/services/creator/types/series"
	videoType "github.com/ystv/web-api/services/creator/types/video"
	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/services/encoder"
)

// Here for validation to ensure we are meeting the interface
var _ creator.BreadcrumbRepo = &Controller{}

// Controller contains our dependency
type Controller struct {
	db     *sqlx.DB
	video  creator.VideoRepo
	series creator.SeriesRepo
}

// NewController creates a new controller
func NewController(db *sqlx.DB, cdn *s3.Client, enc *encoder.Encoder, conf *creator.Config) *Controller {
	return &Controller{db: db, video: video.NewStore(db, cdn, enc, conf)}
}

// Series will return the breadcrumb from SeriesID to root
func (c *Controller) Series(ctx context.Context, seriesID int) ([]breadcrumb.Breadcrumb, error) {
	var s []breadcrumb.Breadcrumb
	err := c.db.SelectContext(ctx, &s,
		`SELECT parent.series_id AS id, parent.url AS url, COALESCE(parent.name, parent.url) AS name, parent.in_url AS use 
		FROM
			video.series node,
			video.series parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND node.series_id = $1
		ORDER BY parent.lft;`, seriesID)
	if err != nil {
		return []breadcrumb.Breadcrumb{}, err
	}
	if len(s) == 0 {
		return []breadcrumb.Breadcrumb{}, series.ErrNotFound
	}
	return s, nil
}

// Video returns the absolute path from a VideoID
func (c *Controller) Video(ctx context.Context, videoID int) ([]breadcrumb.Breadcrumb, error) {
	vB := breadcrumb.Breadcrumb{} // Video breadcrumb
	err := c.db.GetContext(ctx, &vB,
		`SELECT video_id as id, series_id, COALESCE(name, url) as name, url
		FROM video.items
		WHERE video_id = $1`, videoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []breadcrumb.Breadcrumb{}, videoType.ErrNotFound
		}
		return []breadcrumb.Breadcrumb{}, fmt.Errorf("failed to get video breadcrumb: %w", err)
	}
	sB, err := c.Series(ctx, vB.SeriesID)
	if err != nil {
		// Interesting edge-case
		if !errors.Is(err, series.ErrNotFound) {
			return nil, fmt.Errorf("failed to get series breadcrumb: %w", err)
		}
	}
	sB = append(sB, vB)

	return sB, nil
}

// Find will return either a series or a video for a given path
func (c *Controller) Find(ctx context.Context, path string) (breadcrumb.Item, error) {
	s, err := c.series.FromPath(ctx, path)
	if err != nil {
		// Might be a video, so we'll go back a crumb and check for a series
		if errors.Is(err, sql.ErrNoRows) {
			split := strings.Split(path, "/")
			PathWithoutLast := strings.Join(split[:len(split)-1], "/")
			s, err = c.series.FromPath(ctx, PathWithoutLast)
			if err != nil {
				return breadcrumb.Item{}, err
			}
			// Found series
			if len(s.ChildVideos) == 0 {
				// No videos on series
				return breadcrumb.Item{}, videoType.ErrNotFound
			}
			// We've got videos
			for _, v := range s.ChildVideos {
				// Check if video name matches last path
				if v.URL == split[len(split)-1] {
					// Found video
					foundVideo, err := c.video.GetItem(ctx, v.ID)
					if err != nil {
						return breadcrumb.Item{}, fmt.Errorf("failed to get video: %w", err)
					}
					return breadcrumb.Item{Video: foundVideo}, nil
				}
			}
		} else {
			return breadcrumb.Item{}, fmt.Errorf("failed to get series from path: %w", err)
		}
	}
	// Found series
	return breadcrumb.Item{Series: s}, nil
}
