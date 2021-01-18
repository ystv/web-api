package public

import (
	"context"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

// TODO add AND to ensure only public is displayed

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see it's immediate children.
	Series struct {
		SeriesMeta
		ImmediateChildSeries []SeriesMeta `json:"childSeries"`
		ChildVideos          []VideoMeta  `json:"videos"`
	}
	// SeriesMeta is used as a children object for a series
	SeriesMeta struct {
		SeriesID    int         `db:"series_id" json:"id"`
		URL         string      `db:"url" json:"url"`
		SeriesName  string      `db:"name" json:"name"`
		Description string      `db:"description" json:"description"`
		Thumbnail   null.String `db:"thumbnail" json:"thumbnail"`
		Depth       int         `db:"depth" json:"-"`
	}
)

var _ SeriesRepo = &Store{}

// GetSeries provides the immediate children of children and videos
func (m *Store) GetSeries(ctx context.Context, seriesID int) (Series, error) {
	s := Series{}
	s, err := m.GetSeriesMeta(ctx, seriesID)
	if err != nil {
		err = fmt.Errorf("failed to get series meta: %w", err)
		return s, err
	}
	s.ImmediateChildSeries, err = m.GetSeriesImmediateChildrenSeries(ctx, seriesID)
	if err != nil {
		err = fmt.Errorf("failed to get child series: %w", err)
		return s, err
	}
	s.ChildVideos, err = m.VideoOfSeries(ctx, seriesID)
	if err != nil {
		err = fmt.Errorf("failed to get child videos: %w", err)
		return s, err
	}
	return s, nil
}

// GetSeriesMeta provides basic information for only the selected series
// TODO probably want to swap this to return SeriesMeta instead
func (m *Store) GetSeriesMeta(ctx context.Context, seriesID int) (Series, error) {
	s := Series{}
	err := m.db.GetContext(ctx, &s,
		`SELECT series_id, url, name, description, thumbnail
		FROM video.series
		WHERE series_id = $1;`, seriesID)
	return s, err
}

// GetSeriesImmediateChildrenSeries returns series directly below the chosen series
func (m *Store) GetSeriesImmediateChildrenSeries(ctx context.Context, seriesID int) ([]SeriesMeta, error) {
	s := []SeriesMeta{}
	err := m.db.SelectContext(ctx, &s,
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
			where depth = 1;`, seriesID)
	return s, err
}

// GetSeriesFromPath returns a series from a url path
func (m *Store) GetSeriesFromPath(ctx context.Context, path string) (Series, error) {
	s := Series{}
	err := m.db.GetContext(ctx, &s.SeriesID,
		`SELECT series_id
	FROM video.series_paths
	WHERE path = $1
	AND status = 'public'`, path)
	if err != nil {
		return s, err
	}
	s, err = m.GetSeries(ctx, s.SeriesID)
	return s, err
}

// SeriesByYear a virtual series containing child series / videos of content uploaded in that year
func (m *Store) SeriesByYear(ctx context.Context, year int) (Series, error) {
	s := Series{}
	// Putting the child series on pause since it looks like we didn't historically store the
	// the created date of video_boxes, we will need to generate the created_at field at some point
	// based on the child videos upload date
	//
	// err := m.db.SelectContext(ctx, &s.ImmediateChildSeries, `
	// 	SELECT series_id, url, name, description, thumbnail
	// 	FROM video.series
	// 	WHERE EXTRACT(year FROM created_at) = $1;`, year)
	// if err != nil {
	// 	return s, fmt.Errorf("failed to get list of series meta by year: %w", err)
	// }
	err := m.db.SelectContext(ctx, &s.ChildVideos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE EXTRACT(year FROM broadcast_date) = $1;`, year)
	if err != nil {
		return s, fmt.Errorf("failed to get list of video metas by year: %w", err)
	}
	return s, nil
}
