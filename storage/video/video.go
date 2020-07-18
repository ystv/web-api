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
		Duration       time.Duration  `db:"duration" json:"duration"`
		Views          int            `db:"views" json:"views"`
		Tags           pq.StringArray `db:"tags" json:"tags"`
		SeriesPosition null.Int       `db:"series_position" json:"seriesPosition"`
		Status         string         `db:"status" json:"status"`
		Preset         null.String    `db:"preset_name" json:"preset"`
		BroadcastDate  string         `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt      time.Time      `db:"created_at" json:"createdAt"`
		Files          []SQLVideoFile `db:"files" json:"files"`
	}
	// SQLVideoFile represents a more readable VideoFile
	SQLVideoFile struct {
		URI          string `db:"uri" json:"uri"`
		EncodeFormat string `db:"name" json:"encodeFormat"`
		Status       string `db:"status" json:"status"`
		Size         int    `db:"size" json:"size"`
	}
)

type VideoController struct {
	db *sqlx.DB
}

func NewController(db *sqlx.DB) *VideoController {
	return &VideoController{db: db}
}

func (v *VideoController) List(ctx context.Context) error {
	return nil
}

func (v *VideoController) Find(ctx context.Context, id int) error {
	return nil
}

// FindVideoItem returns a VideoItem by it's ID.
func FindVideoItem(ctx context.Context, id int) (*SQLVideoItem, error) {
	v := SQLVideoItem{}
	err := utils.DB.Get(&v,
		`SELECT item.video_id, item.series_id, item.name item_name, item.url,
		item.description, item.thumbnail, EXTRACT(EPOCH FROM item.duration) AS duration,
		item.views, item.tags, item.series_position, item.status,
		preset.name preset_name, trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date, created_at 
		FROM video.items item
		LEFT JOIN video.presets preset ON item.preset = preset.id
		WHERE video_id = $1
		LIMIT 1;`, id)
	log.Print(err)
	if err != nil {
		return nil, err
	}
	err = utils.DB.SelectContext(ctx, &v.Files,
		`SELECT uri, name, status, size
		FROM video.files
		INNER JOIN video.encode_formats ON id = encode_format
		WHERE video_id = $1;`, id)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
