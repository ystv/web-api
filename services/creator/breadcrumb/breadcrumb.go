package breadcrumb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/breadcrumb"
	"github.com/ystv/web-api/services/creator/video"
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
func NewController(db *sqlx.DB, cdn *s3.S3, conf *creator.Config) *Controller {
	return &Controller{db: db, video: video.NewStore(db, cdn, conf)}
}

// Series will return the breadcrumb from SeriesID to root
func (c *Controller) Series(ctx context.Context, seriesID int) (*[]breadcrumb.Breadcrumb, error) {
	s := []breadcrumb.Breadcrumb{}
	// TODO Need a bool to indicate if series is in URL
	err := c.db.SelectContext(ctx, &s,
		`SELECT parent.series_id as id, parent.url as url, COALESCE(parent.name, parent.url) as name
		FROM
			video.series node,
			video.series parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND node.series_id = $1
		ORDER BY parent.lft;`, seriesID)
	if err != nil {
		return nil, err
	}
	return &s, err
}

// Video returns the absolute path from a VideoID
func (c *Controller) Video(ctx context.Context, videoID int) (*[]breadcrumb.Breadcrumb, error) {
	vB := breadcrumb.Breadcrumb{} // Video breadcrumb
	err := c.db.GetContext(ctx, &vB,
		`SELECT video_id as id, series_id, COALESCE(name, url) as name, url
		FROM video.items
		WHERE video_id = $1`, videoID)
	if err != nil {
		err = fmt.Errorf("failed to get video breadcrumb: %w", err)
		return nil, err
	}
	sB, err := c.Series(ctx, vB.SeriesID)
	if err != nil {
		err = fmt.Errorf("failed to get series breadcrumb: %w", err)
		return nil, err
	}
	*sB = append(*sB, vB)

	return sB, err
}

// Find will returns either a series or a video for a given path
func (c *Controller) Find(ctx context.Context, path string) (*breadcrumb.Item, error) {
	s, err := c.series.FromPath(ctx, path)
	if err != nil {
		// Might be a video, so we'll go back a crumb and check for a series
		if err == sql.ErrNoRows {
			split := strings.Split(path, "/")
			PathWithoutLast := strings.Join(split[:len(split)-1], "/")
			s, err := c.series.FromPath(ctx, PathWithoutLast)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// No series, so there will be no videos
					err = fmt.Errorf("No series: %w", err)
					return nil, err
				}
				err = fmt.Errorf("Failed from back a layer: %w", err)
				return nil, err
			}
			// Found series
			if len(*s.ChildVideos) == 0 {
				// No videos on series
				return nil, errors.New("No videos")
			}
			// We've got videos
			for _, v := range *s.ChildVideos {
				// Check if video name matches last path
				if v.URL == split[len(split)-1] {
					// Found video
					foundVideo, err := c.video.GetItem(ctx, v.ID)
					if err != nil {
						return nil, err
					}
					return &breadcrumb.Item{Video: foundVideo}, nil
				}
			}
		} else {
			err = fmt.Errorf("From path %s: %w", path, err)
			return nil, err
		}
	}
	// Found series
	return &breadcrumb.Item{Series: s}, nil
}
