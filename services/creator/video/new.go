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
)

// NewStore returns a new store
func NewStore(db *sqlx.DB, cdn *s3.S3, conf *creator.Config) *Store {
	return &Store{db: db, cdn: cdn, conf: conf}
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
