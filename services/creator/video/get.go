package video

import (
	"context"
	"fmt"

	"github.com/ystv/web-api/services/creator/types/video"
)

// GetItem returns a VideoItem by it's ID.
func (s *Store) GetItem(ctx context.Context, id int) (*video.Item, error) {
	v := video.Item{}
	err := s.db.GetContext(ctx, &v,
		`SELECT item.video_id, item.series_id, item.name video_name, item.url,
		item.description, item.thumbnail, duration,	item.views, item.tags,
		item.status, preset.id preset_id, preset.name preset_name, broadcast_date,
		item.created_at, users.user_id AS created_by_id, users.nickname AS created_by_nick
		FROM video.items item
			LEFT JOIN video.presets preset ON item.preset = preset.id
        	INNER JOIN people.users users ON users.user_id = item.created_by
		WHERE video_id = $1
		LIMIT 1;`, id)
	if err != nil {
		err = fmt.Errorf("failed to get video meta: %w", err)
		return nil, err
	}
	err = s.db.SelectContext(ctx, &v.Files,
		`SELECT uri, name, status, size, mime_type
		FROM video.files
		INNER JOIN video.encode_formats ON id = encode_format
		WHERE video_id = $1;`, id)
	if err != nil {
		err = fmt.Errorf("failed to get video files: %w", err)
		return nil, err
	}
	return &v, nil
}

// ListMeta returns a list of VideoMeta's
func (s *Store) ListMeta(ctx context.Context) (*[]video.Meta, error) {
	v := []video.Meta{}
	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		duration, views, tags, status, broadcast_date,	created_at
		FROM video.items
		ORDER BY broadcast_date DESC;`)
	return &v, err
}

// ListMetaByUser returns a list of VideoMeta's for a given user
func (s *Store) ListMetaByUser(ctx context.Context, userID int) (*[]video.Meta, error) {
	v := []video.Meta{}
	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		duration, views, tags, status, broadcast_date, created_at
		FROM video.items
		WHERE created_by = $1
		ORDER BY broadcast_date DESC;`, userID)
	return &v, err
}

// ListByCalendarMonth returns a list of VideoMeta's for a given month/year
func (s *Store) ListByCalendarMonth(ctx context.Context, year, month int) (*[]video.MetaCal, error) {
	v := []video.MetaCal{}
	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, name, status, broadcast_date
		FROM video.items
		WHERE EXTRACT(YEAR FROM broadcast_date) = $1 AND
		EXTRACT(MONTH FROM broadcast_date) = $2;`, year, month)
	return &v, err
}

// OfSeries returns all the videos belonging to a series
func (s *Store) OfSeries(ctx context.Context, seriesID int) (*[]video.Meta, error) {
	v := []video.Meta{}
	//TODO Update this select to fill all fields
	err := s.db.Select(&v,
		`SELECT video_id, series_id, name video_name, url, broadcast_date,
		views, duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public';`, seriesID)
	return &v, err
}
