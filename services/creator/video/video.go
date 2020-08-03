package video

import (
	"context"
	"log"

	"github.com/lib/pq"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// TODO update schema so duration is not null
// TODO update structs so meta isn't repeated

type (
	// Item represents a more readable VideoItem with
	// an array of associated VideoFiles.
	Item struct {
		Meta
		Files []File `db:"files" json:"files"`
	}
	// File represents a more readable VideoFile.
	File struct {
		URI          string   `db:"uri" json:"uri"`
		EncodeFormat string   `db:"name" json:"encodeFormat"`
		Status       string   `db:"status" json:"status"`
		Size         null.Int `db:"size" json:"size"`
		MimeType     string   `db:"mime_type" json:"mimeType"`
	}
	// TODO make null's pointers, so we can omitempty them during JSON marshal

	// Meta represents just the metadata of a video, used for listing.
	Meta struct {
		ID             int            `db:"video_id" json:"id"`
		SeriesID       int            `db:"series_id" json:"seriesID"`
		Name           string         `db:"video_name" json:"name"`
		URL            string         `db:"url" json:"url"`
		Description    null.String    `db:"description" json:"description,omitempty"`
		Thumbnail      null.String    `db:"thumbnail" json:"thumbnail"`
		Duration       null.Int       `db:"duration" json:"duration"`
		Views          int            `db:"views" json:"views"`
		Tags           pq.StringArray `db:"tags" json:"tags"`
		SeriesPosition null.Int       `db:"series_position" json:"seriesPosition"`
		Status         string         `db:"status" json:"status"`
		Preset         null.String    `db:"preset_name" json:"preset"`
		BroadcastDate  string         `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt      string         `db:"created_at" json:"createdAt"`
	}
	// MetaCal represents simple metadata for a calendar
	MetaCal struct {
		ID            int    `db:"video_id" json:"id"`
		Name          string `db:"name" json:"name"`
		Status        string `db:"status" json:"status"`
		BroadcastDate string `db:"broadcast_date" json:"broadcastDate"`
	}
)

// TODO stop using global DB
// type Controller struct {
// 	db *sqlx.DB
// }

// func NewController(db *sqlx.DB) *Controller {
// 	return &Controller{db: db}
// }

// func (v *Controller) List(ctx context.Context) ([]*Meta, error) {
// 	return nil, nil
// }

// func (v *Controller) Find(ctx context.Context, id int) error {
// 	return nil
// }

// FindVideoItem returns a VideoItem by it's ID.
func FindVideoItem(ctx context.Context, id int) (*Item, error) {
	v := Item{}
	err := utils.DB.GetContext(ctx, &v,
		`SELECT item.video_id, item.series_id, item.name video_name, item.url,
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
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// MetaList returns a list of VideoMeta's
func MetaList(ctx context.Context) (*[]Meta, error) {
	v := []Meta{}
	err := utils.DB.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		EXTRACT(EPOCH FROM duration)::int AS duration, views, tags,
		series_position, status, trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		trim(both '"' from to_json(created_at)::text) AS created_at
		FROM video.items
		ORDER BY broadcast_date DESC;`)
	return &v, err
}

// CalendarList returns a list of VideoMeta's for a given month/year
func CalendarList(ctx context.Context, year int, month int) (*[]MetaCal, error) {
	v := []MetaCal{}
	err := utils.DB.SelectContext(ctx, &v,
		`SELECT video_id, name, status,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date
		FROM video.items
		WHERE EXTRACT(YEAR FROM broadcast_date) = $1 AND
		EXTRACT(MONTH FROM broadcast_date) = $2`, year, month)
	return &v, err
}

// OfSeries returns all the videos belonging to a series
func OfSeries(SeriesID int) ([]Meta, error) {
	v := []Meta{}
	//TODO Update this select to fill all fields
	err := utils.DB.Select(&v,
		`SELECT video_id, series_id, name video_name, url,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public'`, SeriesID)
	if err != nil {
		log.Printf("Failed to select VideoOfSeries: %+v", err)
	}
	return v, err
}
