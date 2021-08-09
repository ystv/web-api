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
	"github.com/ystv/web-api/services/encoder"
	"github.com/ystv/web-api/utils"
)

// NewStore returns a new store
func NewStore(db *sqlx.DB, cdn *s3.S3, enc *encoder.Encoder, conf *creator.Config) *Store {
	return &Store{db: db, cdn: cdn, conf: conf}
}

// NewItem creates a new video item
func (s *Store) NewItem(ctx context.Context, v *video.New) (int, error) {
	// Checking if video file exists
	obj, err := s.cdn.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.conf.IngestBucket),
		Key:    aws.String(v.FileID[:32]),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to find video object \"%s\" in bucket \"%s\": %w", v.FileID[:32], s.conf.IngestBucket, err)
	}

	// Generating timestamp
	v.CreatedAt = time.Now()

	// New video ID, will be filled when created
	var videoID int

	err = utils.Transact(s.db, func(tx *sqlx.Tx) error {
		// Inserting video item record
		itemQuery := `INSERT INTO video.items (series_id, name, url, description, tags,
			status, created_at, created_by, broadcast_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING video_id;`

		err = tx.QueryRowContext(ctx,
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
			Metadata:   obj.Metadata, // TODO: Copy from the soure Content-Type
		})

		if err != nil {
			return fmt.Errorf("failed to copy video object from pending bucket to video bucket: %w", err)
		}

		// Updating DB to reflect this
		fileQuery := `INSERT INTO video.files (video_id, format_id, uri, status, size)
					VALUES ($1, $2, $3, $4, $5);`

		_, err = tx.ExecContext(ctx, fileQuery, videoID, 1, "videos/"+key, "internal", *obj.ContentLength) // TODO make a original encode format
		if err != nil {
			return fmt.Errorf("failed to insert video file row: %w", err)
		}

		return nil
	})
	if err != nil {
		// Since we've wrapped in transaction the DB is safe, will just need to make sure s3 is back to original state
		// TODO: Do we want to care about this outcome?
		s.cdn.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.conf.ServeBucket),
			Key:    aws.String(s.conf.IngestBucket + "/" + v.FileID[:32]),
		})

		return 0, fmt.Errorf("failed to insert create: %w", err)
	}

	// Check if a preset was attached, if so we will start transcoding jobs
	if v.PresetID != 0 {
		err = s.enc.RefreshVideo(ctx, videoID)
		if err != nil {
			return videoID, fmt.Errorf("failed to refresh video: %w", err)
		}
	}
	return videoID, nil
}
