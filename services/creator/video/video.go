package video

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/video"
)

// TODO update schema so duration is not null

var _ creator.VideoRepo = &Store{}

type Store struct {
	db  *sqlx.DB
	cdn *s3.S3
}

func NewStore(db *sqlx.DB, cdn *s3.S3) *Store {
	return &Store{db: db, cdn: cdn}
}

// GetItem returns a VideoItem by it's ID.
func (s *Store) GetItem(ctx context.Context, id int) (*video.Item, error) {
	v := video.Item{}
	err := s.db.GetContext(ctx, &v,
		`SELECT item.video_id, item.series_id, item.name video_name, item.url,
		item.description, item.thumbnail, EXTRACT(EPOCH FROM item.duration)::int AS duration,
		item.views, item.tags, item.status,	preset.id preset_id, preset.name preset_name,
		broadcast_date, item.created_at, users.user_id, users.nickname
		FROM video.items item
			LEFT JOIN video.presets preset ON item.preset = preset.id
        	INNER JOIN people.users users ON users.user_id = item.created_by
		WHERE video_id = $1
		LIMIT 1;`, id)
	log.Print(err)
	if err != nil {
		return nil, err
	}
	err = s.db.SelectContext(ctx, &v.Files,
		`SELECT uri, name, status, size, mime_type
		FROM video.files
		INNER JOIN video.encode_formats ON id = encode_format
		WHERE video_id = $1;`, id)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// ListMeta returns a list of VideoMeta's
func (s *Store) ListMeta(ctx context.Context) (*[]video.Meta, error) {
	v := []video.Meta{}
	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		EXTRACT(EPOCH FROM duration)::int AS duration, views, tags,
		status, trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		trim(both '"' from to_json(created_at)::text) AS created_at
		FROM video.items
		ORDER BY broadcast_date DESC;`)
	return &v, err
}

// ListMetaByUser returns a list of VideoMeta's for a given user
func (s *Store) ListMetaByUser(ctx context.Context, userID int) (*[]video.Meta, error) {
	v := []video.Meta{}
	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, series_id, name video_name, url,
		EXTRACT(EPOCH FROM duration)::int AS duration, views, tags,
		status, trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		trim(both '"' from to_json(created_at)::text) AS created_at
		FROM video.items
		WHERE created_by = $1
		ORDER BY broadcast_date DESC;`, userID)
	return &v, err
}

// ListByCalendarMonth returns a list of VideoMeta's for a given month/year
func (s *Store) ListByCalendarMonth(ctx context.Context, year, month int) (*[]video.MetaCal, error) {
	v := []video.MetaCal{}
	err := s.db.SelectContext(ctx, &v,
		`SELECT video_id, name, status,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date
		FROM video.items
		WHERE EXTRACT(YEAR FROM broadcast_date) = $1 AND
		EXTRACT(MONTH FROM broadcast_date) = $2`, year, month)
	return &v, err
}

// OfSeries returns all the videos belonging to a series
func (s *Store) OfSeries(ctx context.Context, seriesID int) (*[]video.Meta, error) {
	v := []video.Meta{}
	//TODO Update this select to fill all fields
	err := s.db.Select(&v,
		`SELECT video_id, series_id, name video_name, url,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public'`, seriesID)
	if err != nil {
		log.Printf("Failed to select VideoOfSeries: %+v", err)
	}
	return &v, err
}

// NewItem creates a new video item
func (s *Store) NewItem(ctx context.Context, v *video.NewVideo) error {
	// Checking if video file exists
	obj, err := s.cdn.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String("pending"),
		Key:    aws.String(v.FileID[:32]),
	})
	if err != nil {
		log.Printf("NewItem object find fail: %v", err)
		return err
	}

	// Generating timestamp
	v.CreatedAt = time.Now()

	// Inserting video item record
	itemQuery := `INSERT INTO video.items (series_id, name, url, description, tags,
		status, created_at, created_by, broadcast_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING video_id;`
	var videoID int

	err = s.db.QueryRowContext(ctx,
		itemQuery, &v.SeriesID, &v.Name, &v.URLName, &v.Description, pq.Array(v.Tags), &v.PublishType, &v.CreatedAt, &v.CreatedBy, &v.BroadcastDate).Scan(&videoID)
	if err != nil {
		log.Printf("NewItem failed to insert: %v", err)
		return err
	}
	extension := strings.Split(*obj.Metadata["Filename"], ".")
	key := fmt.Sprintf("%d_%d_%s_%s.%s", v.BroadcastDate.Year(), videoID, v.URLName, getSeason(v.BroadcastDate), extension[1])

	// Copy from pending bucket to main video bucket
	_, err = s.cdn.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String("videos"),
		CopySource: aws.String("pending/" + v.FileID[:32]),
		Key:        aws.String(key),
	})

	// Updating DB to reflect this
	fileQuery := `INSERT INTO video.files (video_id, uri, status, encode_format, size)
				VALUES ($1, $2, $3, $4, $5);`

	_, err = s.db.ExecContext(ctx, fileQuery, videoID, "videos/"+key, "internal", 1, *obj.ContentLength) // TODO make a original encode format
	if err != nil {
		log.Printf("NewItem failed to insert video file: %v", err)
		return err
	}

	return nil
}

func getSeason(t time.Time) string {
	m := int(t.Month())
	switch {
	case m >= 9 && m <= 12:
		return "aut"
	case m >= 1 && m <= 6:
		return "spr"
	default:
		return "sum"
	}
}
