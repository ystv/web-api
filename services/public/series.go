package public

import (
	"context"
	"fmt"

	"github.com/ystv/web-api/utils"
)

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see its immediate children.
	Series struct {
		SeriesMeta
		ImmediateChildSeries []SeriesMeta `json:"childSeries"`
		ChildVideos          []VideoMeta  `json:"videos"`
	}
	// SeriesMeta is used as a children object for a series
	SeriesMeta struct {
		SeriesID    int    `db:"series_id" json:"id"`
		URL         string `db:"url" json:"url"`
		SeriesName  string `db:"name" json:"name"`
		Description string `db:"description" json:"description"`
		Thumbnail   string `db:"thumbnail" json:"thumbnail"`
		Depth       int    `db:"depth" json:"-"`
	}
)

// GetSeries provides the immediate children of children and videos
func (s *Store) GetSeries(ctx context.Context, seriesID int) (Series, error) {
	series, err := s.GetSeriesFullMeta(ctx, seriesID)
	if err != nil {
		return series, fmt.Errorf("failed to get series meta: %w", err)
	}

	series.ImmediateChildSeries, err = s.GetSeriesImmediateChildrenSeries(ctx, seriesID)
	if err != nil {
		return series, fmt.Errorf("failed to get child series: %w", err)
	}

	series.ChildVideos, err = s.VideoOfSeries(ctx, seriesID)
	if err != nil {
		return series, fmt.Errorf("failed to get child videos: %w", err)
	}

	return series, nil
}

// GetSeriesMeta provides basic information for only the selected series
func (s *Store) GetSeriesMeta(ctx context.Context, seriesID int) (SeriesMeta, error) {
	var series SeriesMeta
	//nolint:musttag
	err := s.db.GetContext(ctx, &series,
		`SELECT series_id, url, name, description, thumbnail
		FROM video.series
		WHERE series_id = $1
		AND status = 'public';`, seriesID)

	return series, err
}

// GetSeriesFullMeta provides basic information for only the selected series
func (s *Store) GetSeriesFullMeta(ctx context.Context, seriesID int) (Series, error) {
	var series Series
	//nolint:musttag
	err := s.db.GetContext(ctx, &series,
		`SELECT series_id, url, name, description, thumbnail
		FROM video.series
		WHERE series_id = $1
		AND status = 'public';`, seriesID)
	if err != nil {
		return Series{}, fmt.Errorf("failed to get series: %w", err)
	}

	err = s.db.SelectContext(ctx, &series.ChildVideos,
		`SELECT video_id, series_id, name, url, description, thumbnail, broadcast_date, views, duration
		FROM video.items
		WHERE series_id = $1
		AND status = 'public';`, seriesID)
	if err != nil {
		return Series{}, fmt.Errorf("failed to get child videos: %w", err)
	}

	return series, nil
}

// GetSeriesImmediateChildrenSeries returns series directly below the chosen series
func (s *Store) GetSeriesImmediateChildrenSeries(ctx context.Context, seriesID int) ([]SeriesMeta, error) {
	var seriesMeta []SeriesMeta

	err := s.db.SelectContext(ctx, &seriesMeta,
		`SELECT series_id, url, name, description, thumbnail from (
			SELECT 
						node.series_id, node.url, node.name, node.description, node.thumbnail, node.status,
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
								node.lft between parent.lft AND parent.rgt
								AND node.series_id = $1
							GROUP BY node.series_id
							ORDER BY node.lft
						) AS sub_tree
					WHERE
						node.lft BETWEEN parent.lft AND parent.rgt
						AND node.lft BETWEEN sub_parent.lft AND sub_parent.rgt
						AND sub_parent.series_id = sub_tree.series_id
					GROUP BY node.series_id, sub_tree.depth, node.lft
					ORDER BY node.lft
			) AS queries
			WHERE depth = 1 AND
			status = 'public';`, seriesID)

	return utils.NonNil(seriesMeta), err
}

// GetSeriesFromPath returns a series from an url path
func (s *Store) GetSeriesFromPath(ctx context.Context, path string) (Series, error) {
	var series Series

	err := s.db.GetContext(ctx, &series.SeriesID,
		`SELECT series_id
	FROM video.series_paths
	WHERE path = $1
	AND status = 'public';`, path)
	if err != nil {
		return series, err
	}

	series, err = s.GetSeries(ctx, series.SeriesID)
	return series, err
}

// GetSeriesByYear a virtual series containing child series / videos of content uploaded in that year
func (s *Store) GetSeriesByYear(ctx context.Context, year int) (Series, error) {
	var series Series
	// Putting the child series on pause since it looks like we didn't historically store the
	// created date of video_boxes, we will need to generate the created_at field at some point
	// based on the child videos upload date
	//
	// err := m.db.SelectContext(ctx, &s.ImmediateChildSeries, `
	// 	SELECT series_id, url, name, description, thumbnail
	// 	FROM video.series
	// 	WHERE EXTRACT(year FROM created_at) = $1;`, year)
	// if err != nil {
	// 	return s, fmt.Errorf("failed to get list of series meta by year: %w", err)
	// }
	err := s.db.SelectContext(ctx, &series.ChildVideos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.items
		WHERE EXTRACT(year FROM broadcast_date) = $1 AND
		status = 'public';`, year)
	if err != nil {
		return series, fmt.Errorf("failed to get list of video metas by year: %w", err)
	}

	return series, nil
}

// Search performs a full-text search on video library
//
// Uses postgres full-text search, video and series tables to try to make some sense
func (s *Store) Search(ctx context.Context, query string) (Series, error) {
	var series Series

	err := s.db.SelectContext(ctx, &series.ChildVideos,
		`SELECT
			video_id,
			series_id,
			name,
			url,
			description,
			thumbnail,
		 	broadcast_date,
		 	views,
			duration
		FROM (
			SELECT
		  		video.video_id,
		  		video.series_id,
		  		video.name,
		  		video.url,
		  		video.description,
		  		video.thumbnail,
		  		video.views,
		  		video.duration,
		  		video.tags,
		  		video.broadcast_date,
				
				to_tsvector(video.name) || ' ' ||
				to_tsvector(video.description) || ' ' ||
				to_tsvector(unnest(video.tags)) || ' ' ||
				to_tsvector(CAST(video.broadcast_date AS text)) || ' ' ||
				to_tsvector(unnest(array_agg(series.name))) || ' ' ||
				to_tsvector(unnest(array_agg(series.description)))
		  		AS document
			FROM video.items video
			INNER JOIN video.series series ON video.series_id = series.series_id
			WHERE video.status = 'public'
			GROUP BY video.video_id) p_search,

			ts_rank_cd(p_search.document, replace(plainto_tsquery($1)::text, '&', '|')::tsquery) rank

	WHERE p_search.document @@ replace(plainto_tsquery($1)::text, '&', '|')::tsquery
	ORDER BY rank DESC;`, query)
	if err != nil {
		return Series{}, fmt.Errorf("failed to search videos: %w", err)
	}

	return series, nil
}
