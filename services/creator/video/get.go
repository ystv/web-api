package video

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/types/video"
)

// GetItem returns a VideoItem by its ID.
func (s *Store) GetItem(ctx context.Context, videoID int) (video.Item, error) {
	v := video.Item{}
	err := s.db.GetContext(ctx, &v,
		`SELECT item.video_id, item.series_id, item.name video_name, item.url,
		item.description, item.thumbnail, duration,	item.views, item.tags,
		item.status, preset.preset_id, preset.name preset_name, broadcast_date,
		item.created_at, users.user_id AS created_by_id, users.nickname AS created_by_nick
		FROM video.items item
			LEFT JOIN video.encode_presets preset ON item.preset_id = preset.preset_id
        	INNER JOIN people.users users ON users.user_id = item.created_by
		WHERE video_id = $1
		LIMIT 1;`, videoID)
	if err != nil {
		err = fmt.Errorf("failed to get video meta: %w", err)
		return video.Item{}, err
	}

	err = s.db.SelectContext(ctx, &v.Files,
		`SELECT uri, name, status, size, mime_type
		FROM video.files file
		INNER JOIN video.encode_formats format ON file.format_id = format.format_id
		WHERE video_id = $1;`, videoID)
	if err != nil {
		err = fmt.Errorf("failed to get video files: %w", err)
		return video.Item{}, err
	}

	return v, nil
}

// ListMeta returns a list of VideoMeta's
func (s *Store) ListMeta(ctx context.Context) ([]video.Meta, error) {
	var v []video.Meta

	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		duration, views, tags, status, broadcast_date,	created_at
		FROM video.items
		ORDER BY broadcast_date DESC;`)

	return v, err
}

// ListMetaByUser returns a list of VideoMeta's for a given user
func (s *Store) ListMetaByUser(ctx context.Context, userID int) ([]video.Meta, error) {
	var v []video.Meta

	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		duration, views, tags, status, broadcast_date, created_at
		FROM video.items
		WHERE created_by = $1
		ORDER BY broadcast_date DESC;`, userID)

	return v, err
}

// ListByCalendarMonth returns a list of VideoMeta's for a given month/year
func (s *Store) ListByCalendarMonth(ctx context.Context, year, month int) ([]video.MetaCal, error) {
	var v []video.MetaCal

	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, name, status, broadcast_date
		FROM video.items
		WHERE EXTRACT(YEAR FROM broadcast_date) = $1 AND
		EXTRACT(MONTH FROM broadcast_date) = $2;`, year, month)

	return v, err
}

// OfSeries returns all the videos belonging to a series
func (s *Store) OfSeries(ctx context.Context, seriesID int) ([]video.Meta, error) {
	var v []video.Meta

	//TODO Update this select to fill all fields
	err := s.db.Select(&v,
		`SELECT video_id, series_id, name video_name, url,
		duration, views, tags, status, broadcast_date, created_at
		FROM video.items
		WHERE series_id = $1;`, seriesID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []video.Meta{}, series.ErrChildrenVideosNotFound
		}
		return []video.Meta{}, err
	}

	return v, nil
}

// Search performs a full-text search on video library
//
// Uses postgres full-text search, video and series tables to try to make some sense
func (s *Store) Search(ctx context.Context, query string) ([]video.Meta, error) {
	var videos []video.Meta

	err := s.db.SelectContext(ctx, &videos,
		`SELECT
			video_id,
			name,
			url,
			description,
			thumbnail,
			broadcast_date,
			views,
			duration,
			tags,
			status
   		FROM (
			SELECT
				video.video_id,
				video.name,
				video.url,
				video.description,
				video.thumbnail,
				video.views,
				video.duration,
				video.tags,
				video.status,
				video.broadcast_date broadcast_date,

				to_tsvector(video.name) || ' ' ||
				to_tsvector(video.description) || ' ' ||
				to_tsvector(unnest(video.tags)) || ' ' ||
				to_tsvector(CAST(video.broadcast_date AS text)) || ' ' ||
				to_tsvector(unnest(array_agg(series.name))) || ' ' ||
				to_tsvector(unnest(array_agg(series.description)))
				AS document
   			FROM video.items video
			INNER JOIN video.series series ON video.series_id = series.series_id
			GROUP BY video.video_id) p_search,
			
			ts_rank_cd(p_search.document, replace(plainto_tsquery($1)::text, '&', '|')::tsquery) rank

   WHERE p_search.document @@ replace(plainto_tsquery($1)::text, '&', '|')::tsquery
   ORDER BY rank DESC;`, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}

	return videos, nil
}
