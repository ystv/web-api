package video

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/video"
	"github.com/ystv/web-api/utils"
)

// TODO update schema so duration is not null

var _ creator.VideoRepo = &Store{}

// Store encapsulates our dependencies
type Store struct {
	db   *sqlx.DB
	cdn  *s3.S3
	conf *creator.Config
}

// NewStore returns a new store
func NewStore(db *sqlx.DB, cdn *s3.S3, conf *creator.Config) *Store {
	return &Store{db: db, cdn: cdn, conf: conf}
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
		EXTRACT(MONTH FROM broadcast_date) = $2;`, year, month)
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
		WHERE series_id = $1 AND status = 'public';`, seriesID)
	return &v, err
}

// NewItem creates a new video item
// TODO I think this needs to be redesigned more like a transaction so we can safely fail anywhere.
// TODO return new video ID
func (s *Store) NewItem(ctx context.Context, v *video.New) error {
	// Checking if video file exists
	obj, err := s.cdn.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.conf.IngestBucket),
		Key:    aws.String(v.FileID[:32]),
	})
	if err != nil {
		err = fmt.Errorf("failed to find video object in s3: %w", err)
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
		err = fmt.Errorf("failed to insert video item: %w", err)
		return err
	}
	extension := strings.Split(*obj.Metadata["Filename"], ".")
	key := fmt.Sprintf("%d_%d_%s_%s.%s", v.BroadcastDate.Year(), videoID, v.URLName, getSeason(v.BroadcastDate), extension[1])

	// Copy from pending bucket to main video bucket
	_, err = s.cdn.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.conf.ServeBucket),
		CopySource: aws.String(s.conf.IngestBucket + "/" + v.FileID[:32]),
		Key:        aws.String(key),
	})
	if err != nil {
		err = fmt.Errorf("failed to copy video object from pending bucket to video bucket: %w", err)
		return err
	}

	// Updating DB to reflect this
	fileQuery := `INSERT INTO video.files (video_id, uri, status, encode_format, size)
				VALUES ($1, $2, $3, $4, $5);`

	_, err = s.db.ExecContext(ctx, fileQuery, videoID, "videos/"+key, "internal", 1, *obj.ContentLength) // TODO make a original encode format
	if err != nil {
		err = fmt.Errorf("failed to insert video file row: %w", err)
		return err
	}

	return nil
}

// DeleteItem Removes a video. The video will still be present in the database, files
// and visible to users with high enough access
func (s *Store) DeleteItem(ctx context.Context, videoID, userID int) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE video.items SET
			deleted_at = NOW()
			deleted_by = $2
		WHERE video_id = $1;`, videoID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete video item: %w", err)
	}
	return nil
}

// DeleteItemPermanently removes a video entirely, including the associated video files
func (s *Store) DeleteItemPermanently(ctx context.Context, videoID int) error {
	// To delete a video we will need to delete all child objects database first
	// * Video hits
	// * Video files
	// Then we will need to delete the object files
	// * VOD files
	// * Original master
	fileURLs := []string{}
	// Wrapped in transaction so we can rollback if it fails, however
	// S3 doesn't support transactions so only database is protected
	err := utils.Transact(s.db, func(tx *sqlx.Tx) error {
		// Get child files
		err := tx.SelectContext(ctx, &fileURLs, `
			SELECT  uri
			FROM video.files
			WHERE video_id = $1;`, videoID)
		if err != nil {
			return fmt.Errorf("failed to find video file URLs: %w", err)
		}

		// First delete from database
		_, err = tx.ExecContext(ctx, `DELETE FROM video.files WHERE video_id = $1;`, videoID)
		if err != nil {
			return fmt.Errorf("failed to delete video file from database: %w", err)
		}

		// Then deleting from object store
		for _, file := range fileURLs {
			_, err := s.cdn.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(s.conf.ServeBucket),
				Key:    aws.String(file),
			})
			if err != nil {
				return fmt.Errorf("failed to delete video file object: %w", err)
			}

			// Finally removing the video item / meta from the database
			_, err = tx.ExecContext(ctx, `DELETE FROM video.items WHERE video_id = $1`, videoID)
			if err != nil {
				return fmt.Errorf("failed to delete video item from database: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to permanently delete video \"%d\": %w", videoID, err)
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
