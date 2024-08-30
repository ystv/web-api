package series

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/services/encoder"
)

// Here for validation to ensure we are meeting the interface
var _ creator.SeriesRepo = &Controller{}

// Controller contains our dependencies
type Controller struct {
	db    *sqlx.DB
	video creator.VideoRepo
}

// NewController creates a new controller
func NewController(db *sqlx.DB, cdn *s3.S3, enc *encoder.Encoder, conf *creator.Config) *Controller {
	return &Controller{db: db, video: video.NewStore(db, cdn, enc, conf)}
}

// Get provides the immediate children of series and videos
func (c *Controller) Get(ctx context.Context, seriesID int) (series.Series, error) {
	s := series.Series{}
	meta, err := c.GetMeta(ctx, seriesID)
	if err != nil {
		if errors.Is(err, series.ErrMetaNotFound) {
			return series.Series{}, series.ErrNotFound
		}
		return series.Series{}, fmt.Errorf("failed to get series meta: %w", err)
	}
	s.Meta = meta

	// Allowing these children not found errors to be ignored since they are optional

	s.ImmediateChildSeries, err = c.ImmediateChildrenSeries(ctx, seriesID)
	if err != nil {
		if !errors.Is(err, series.ErrChildrenSeriesNotFound) {
			return series.Series{}, fmt.Errorf("failed to get child series: %w", err)
		}
	}
	s.ChildVideos, err = c.video.OfSeries(ctx, seriesID)
	if err != nil {
		if !errors.Is(err, series.ErrChildrenVideosNotFound) {
			return series.Series{}, fmt.Errorf("failed to get child videos: %w", err)
		}
	}
	return s, nil
}

// GetMeta provides basic information for only the selected series
func (c *Controller) GetMeta(ctx context.Context, seriesID int) (series.Meta, error) {
	s := series.Meta{}
	err := c.db.GetContext(ctx, &s,
		`SELECT series_id, url, name, description, thumbnail
		FROM video.series
		WHERE series_id = $1`, seriesID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return series.Meta{}, series.ErrMetaNotFound
		}
		return series.Meta{}, err
	}
	return s, nil
}

// ImmediateChildrenSeries returns series directly below the chosen series
func (c *Controller) ImmediateChildrenSeries(ctx context.Context, SeriesID int) ([]series.Meta, error) {
	var s []series.Meta
	err := c.db.SelectContext(ctx, &s,
		`SELECT * from (
			SELECT 
						node.series_id, node.url, node.name, node.description, node.thumbnail,
						(COUNT(parent.*) - (sub_tree.depth + 1)) AS depth
					FROM
						video.series AS node,
						video.series AS parent,
						video.series AS sub_parent,
						(
							SELECT node.series_id, (COUNT(parent.*) - 1) AS depth
							FROM
								video.series AS node,
								video.series AS parent
							WHERE
								node.lft between parent.lft and parent.rgt
								and node.series_id = $1
							GROUP BY node.series_id
							ORDER BY node.lft ASC
						) AS sub_tree
					WHERE
						node.lft BETWEEN parent.lft AND parent.rgt
						AND node.lft BETWEEN sub_parent.lft AND sub_parent.rgt
						AND sub_parent.series_id = sub_tree.series_id
					GROUP BY node.series_id, sub_tree.depth
					ORDER BY node.lft asc
			) as queries
			where depth = 1;`, SeriesID)
	if err != nil {
		return []series.Meta{}, err
	}
	if len(s) == 0 {
		return []series.Meta{}, series.ErrChildrenSeriesNotFound
	}
	return s, nil
}

// List returns all series in the DB including their depth
func (c *Controller) List(ctx context.Context) ([]series.Meta, error) {
	var s []series.Meta
	err := c.db.SelectContext(ctx, s,
		`SELECT
			child.series_id, child.url, child.name, child.description, child.thumbnail,
			(COUNT(parent.*) -1) AS depth
		FROM
			video.series child,
			video.series parent
		WHERE
			child.lft BETWEEN parent.lft AND parent.rgt
		GROUP BY child.series_id
		ORDER BY child.lft ASC;`)
	if err != nil {
		return []series.Meta{}, err
	}
	if len(s) == 0 {
		return []series.Meta{}, series.ErrNotFound
	}
	return s, nil
}

// AllBelow returns all series below a certain series including depth
func (c *Controller) AllBelow(ctx context.Context, SeriesID int) ([]series.Meta, error) {
	var s []series.Meta
	err := c.db.SelectContext(ctx, &s,
		`SELECT 
			node.series_id, node.url node.name, node.description, node.thumbnail,
			(COUNT(parent.*) - (sub_tree.depth + 1)) AS depth
		FROM
			video.series AS node,
			video.series AS parent,
			video.series AS sub_parent,
			(
				SELECT node.series_id, (COUNT(parent.*) - 1) AS depth
				FROM
					video.series AS node,
					video.series AS parent
				WHERE
					node.lft between parent.lft and parent.rgt
					and node.series_id = $1
				GROUP BY node.series_id
				ORDER BY node.lft ASC
			) AS sub_tree
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND node.lft BETWEEN sub_parent.lft AND sub_parent.rgt
			AND sub_parent.series_id = sub_tree.series_id
		GROUP BY node.series_id, sub_tree.depth
		ORDER BY node.lft ASC;`, SeriesID)
	if err != nil {
		return []series.Meta{}, err
	}
	if len(s) == 0 {
		return []series.Meta{}, series.ErrNotFound
	}
	return s, nil
}

// FromPath will return a series from a given path
func (c *Controller) FromPath(ctx context.Context, path string) (series.Series, error) {
	s := series.Series{}
	err := c.db.GetContext(ctx, s.SeriesID, `SELECT series_id FROM video.series_paths WHERE path = $1`, path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return series.Series{}, series.ErrNotFound
		}
		return series.Series{}, fmt.Errorf("failed to get series from path: %w", err)
	}
	s, err = c.Get(ctx, s.SeriesID)
	if err != nil {
		err = fmt.Errorf("failed to get series data: %w", err)
		return series.Series{}, err
	}
	return s, err
}
