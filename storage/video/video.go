package video

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// TODO update schema so duration is not null

type (
	// SQLVideoItem represents a more readable VideoItem with
	// an array of associated VideoFiles.
	SQLVideoItem struct {
		ID             int            `db:"video_id" json:"videoID"`
		SeriesID       int            `db:"series_id" json:"seriesID"`
		Name           string         `db:"item_name" json:"name"`
		URL            string         `db:"url" json:"url"`
		Description    null.String    `db:"description" json:"description"`
		Thumbnail      null.String    `db:"thumbnail" json:"thumbnail"`
		Duration       null.Int       `db:"duration" json:"duration"`
		Views          int            `db:"views" json:"views"`
		Tags           pq.StringArray `db:"tags" json:"tags"`
		SeriesPosition null.Int       `db:"series_position" json:"seriesPosition"`
		Status         string         `db:"status" json:"status"`
		Preset         null.String    `db:"preset_name" json:"preset"`
		BroadcastDate  string         `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt      time.Time      `db:"created_at" json:"createdAt"`
		Files          []SQLVideoFile `db:"files" json:"files"`
	}
	// SQLVideoFile represents a more readable VideoFile.
	SQLVideoFile struct {
		URI          string   `db:"uri" json:"uri"`
		EncodeFormat string   `db:"name" json:"encodeFormat"`
		Status       string   `db:"status" json:"status"`
		Size         null.Int `db:"size" json:"size"`
		MimeType     string   `db:"mime_type" json:"mimeType"`
	}
	// SQLVideoMeta represents just the metadata of a video, used for listing.
	SQLVideoMeta struct {
		ID             int            `db:"video_id" json:"videoID"`
		SeriesID       int            `db:"series_id" json:"seriesID"`
		Name           string         `db:"name" json:"name"`
		URL            string         `db:"url" json:"url"`
		Duration       null.Int       `db:"duration" json:"duration"`
		Views          int            `db:"views" json:"views"`
		Tags           pq.StringArray `db:"tags" json:"tags"`
		SeriesPosition null.Int       `db:"series_position" json:"seriesPosition"`
		Status         string         `db:"status" json:"status"`
		BroadcastDate  string         `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt      string         `db:"created_at" json:"createdAt"`
	}
	// SQLVideoMetaCal represents simple metadata for a calendar
	SQLVideoMetaCal struct {
		ID            int    `db:"video_id" json:"videoID"`
		Name          string `db:"name" json:"name"`
		BroadcastDate string `db:"broadcast_date" json:"broadcastDate"`
	}
)

type Controller struct {
	db *sqlx.DB
}

func NewController(db *sqlx.DB) *Controller {
	return &Controller{db: db}
}

func (v *Controller) List(ctx context.Context) ([]*SQLVideoMeta, error) {
	return nil, nil
}

func (v *Controller) Find(ctx context.Context, id int) error {
	return nil
}

// FindVideoItem returns a VideoItem by it's ID.
func FindVideoItem(ctx context.Context, id int) (*SQLVideoItem, error) {
	v := SQLVideoItem{}
	err := utils.DB.Get(&v,
		`SELECT item.video_id, item.series_id, item.name item_name, item.url,
		item.description, item.thumbnail, EXTRACT(EPOCH FROM item.duration)::int AS duration,
		item.views, item.tags, item.series_position, item.status,
		preset.name preset_name, trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date, created_at 
		FROM video.items item
		LEFT JOIN video.presets preset ON item.preset = preset.id
		WHERE video_id = $1
		LIMIT 1;`, id)
	if err != nil {
		return nil, err
	}
	err = utils.DB.SelectContext(ctx, &v.Files,
		`SELECT uri, name, status, size, mime_type
		FROM video.files
		INNER JOIN video.encode_formats ON id = encode_format
		WHERE video_id = $1;`, id)
	log.Print(err)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// MetaList returns a list of VideoMeta's
func MetaList(ctx context.Context) (*[]SQLVideoMeta, error) {
	v := []SQLVideoMeta{}
	err := utils.DB.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name, url,
		EXTRACT(EPOCH FROM duration)::int AS duration, views, tags,
		series_position, status, trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		trim(both '"' from to_json(created_at)::text) AS created_at
		FROM video.items
		ORDER BY video_id;`)
	return &v, err
}

// CalendarList returns a list of VideoMeta's for a given month/year
func CalendarList(ctx context.Context, year int, month int) (*[]SQLVideoMetaCal, error) {
	v := []SQLVideoMetaCal{}
	err := utils.DB.SelectContext(ctx, &v,
		`SELECT video_id, name,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date
		FROM video.items
		WHERE EXTRACT(YEAR FROM broadcast_date) = $1 AND
		EXTRACT(MONTH FROM broadcast_date) = $2`, year, month)
	log.Print(v)
	return &v, err
}
