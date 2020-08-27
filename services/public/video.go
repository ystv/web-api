package public

import (
	"context"
	"fmt"

	_ "github.com/lib/pq" // for DB, although likely not needed
	"gopkg.in/guregu/null.v4"
)

type (
	// VideoItem represents public info about video item.
	VideoItem struct {
		VideoMeta
		Files []VideoFile `json:"files"`
	}
	// VideoFile represents each file that a video item has stored.
	VideoFile struct {
		URI      string `json:"uri"`
		MimeType string `db:"mime_type" json:"mimeType"`
		Mode     string `db:"mode" json:"mode"`
		Width    int    `db:"width" json:"width"`
		Height   int    `db:"height" json:"height"`
	}
	// VideoMeta represents basic information about the videoitem used for listing.
	VideoMeta struct {
		VideoID       int         `db:"video_id" json:"id"`
		SeriesID      int         `db:"series_id" json:"seriesID"`
		Name          string      `db:"name" json:"name"`
		URL           string      `db:"url" json:"url"`
		Description   string      `db:"description" json:"description"`
		Thumbnail     null.String `db:"thumbnail" json:"thumbnail"`
		BroadcastDate string      `db:"broadcast_date" json:"broadcastDate"`
		Views         int         `db:"views" json:"views"`
		Duration      null.Int    `db:"duration" json:"duration"`
	}
)

var _ VideoRepo = &Store{}

// ListVideo returns all video metadata
func (m *Store) ListVideo(ctx context.Context, offset int, page int) (*[]VideoMeta, error) {
	v := []VideoMeta{}
	// TODO Change pagination method
	err := m.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE status = 'public'
		ORDER BY broadcast_date DESC
		OFFSET $1 LIMIT $2;`, page, offset)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// GetVideo returns a VideoItem, including the files, based on a given VideoItem ID.
func (m *Store) GetVideo(ctx context.Context, videoID int) (*VideoItem, error) {
	v := VideoItem{}
	err := m.db.GetContext(ctx, &v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
	views, EXTRACT(EPOCH FROM duration)::int AS duration,
	trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date
	FROM video.items
	WHERE video_id = $1
	AND status = 'public'
	LIMIT 1;`, videoID)
	if err != nil {
		err = fmt.Errorf("failed to get video meta: %w", err)
		return nil, err
	}
	err = m.db.SelectContext(ctx, &v.Files,
		`SELECT uri, mime_type, mode, width, height
	FROM video.files
	INNER JOIN video.encode_formats ON id = encode_format
	WHERE status = 'public'
	AND video_id = $1`, videoID)
	if err != nil {
		err = fmt.Errorf("failed to get video files: %w", err)
		return nil, err
	}
	return &v, nil
}

// VideoOfSeries returns all the videos belonging to a series
func (m *Store) VideoOfSeries(ctx context.Context, seriesID int) ([]VideoMeta, error) {
	v := []VideoMeta{}
	err := m.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public'
		ORDER BY series_position ASC;`, seriesID)
	if err != nil {
		return nil, err
	}
	return v, err
}
