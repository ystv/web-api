package public

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq" // for DB, although likely not needed
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
		VideoID       int       `db:"video_id" json:"id"`
		SeriesID      int       `db:"series_id" json:"seriesID"`
		Name          string    `db:"name" json:"name"`
		URL           string    `db:"url" json:"url"`
		Description   string    `db:"description" json:"description"`
		Thumbnail     string    `db:"thumbnail" json:"thumbnail"`
		BroadcastDate time.Time `db:"broadcast_date" json:"broadcastDate"`
		Views         int       `db:"views" json:"views"`
		Duration      int       `db:"duration" json:"duration"`
	}
)

var _ VideoRepo = &Store{}

// ListVideo returns all video metadata
func (m *Store) ListVideo(ctx context.Context, offset int, page int) (*[]VideoMeta, error) {
	v := []VideoMeta{}
	// TODO Change pagination method
	// TODO Do a double check on if we need to convert broadcast date
	err := m.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
		broadcast_date,	views, duration
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
	views, duration, broadcast_date
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
		broadcast_date,	views, duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public'
		ORDER BY series_position ASC;`, seriesID)
	if err != nil {
		return nil, err
	}
	return v, err
}

// Search performs a full-text search on video library
//
// Uses postgres' full-text search, video and series tables to try to make some sense
func (m *Store) Search(ctx context.Context, search string) ([]VideoMeta, error) {
	v := []VideoMeta{}
	res, err := m.db.QueryContext(ctx,
		`SELECT id, title, description, tags, broadcast_date
	FROM (SELECT	video.video_id id,
		  video.name title,
		  video.description description,
		  video.tags tags,
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
	GROUP BY video.video_id) p_search
	WHERE p_search.document @@ replace(plainto_tsquery('$1')::text, '&', '|')::tsquery;`, search)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}
	err = res.Scan(&v)
	if err != nil {
		return nil, fmt.Errorf("failed to scan search struct: %w", err)
	}
	return v, nil
}
