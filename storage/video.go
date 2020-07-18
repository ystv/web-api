package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/ystv/web-api/utils"
)

type (
	SQLVideoItem struct {
		ID             int            `db:"video_id"`
		SeriesID       int            `db:"series_id"`
		Name           string         `db:"item_name"`
		URL            string         `db:"url"`
		Description    sql.NullString `db:"description"`
		Thumbnail      sql.NullString `db:"thumbnail"`
		Duration       time.Duration  `db:"duration"`
		Views          int            `db:"views"`
		Tags           pq.StringArray `db:"tags"`
		SeriesPosition sql.NullInt32  `db:"series_position"`
		Status         string         `db:"status"`
		Preset         sql.NullString `db:"preset_name"`
		BroadcastDate  time.Time      `db:"broadcast_date"`
		CreatedAt      time.Time      `db:"created_at"`
		Files          []SQLVideoFile `db:"files"`
	}
	SQLVideoFile struct {
		URI          string `db:"uri"`
		EncodeFormat string `db:"name"`
		Status       string `db:"status"`
		Size         int    `db:"size"`
	}
)

// VideoItem represents a more readable VideoItem with array of
// associated VideoFiles.
func VideoItem(ctx context.Context, id int) (*SQLVideoItem, error) {
	v := SQLVideoItem{}
	err := utils.DB.Get(&v,
		`SELECT item.video_id, item.series_id, item.name item_name, item.url,
		item.description, item.thumbnail, EXTRACT(EPOCH FROM item.duration) AS duration,
		item.views, item.tags, item.series_position, preset.name preset_name,
		broadcast_date, created_at 
		FROM video.items item
		LEFT JOIN video.presets preset ON item.preset = preset.id
		WHERE video_id = $1
		LIMIT 1;`, id)
	log.Printf("Error1: %+v", err)
	if err != nil {
		return nil, err
	}
	err = utils.DB.SelectContext(ctx, &v.Files,
		`SELECT uri, name, status, size
		FROM video.files
		INNER JOIN video.encode_formats ON id = encode_format
		WHERE video_id = $1;`, id)
	log.Printf("Error2: %+v", err)
	if err != nil {
		return nil, err
	}
	log.Printf("Video: %+v", v)
	return &v, nil
}
