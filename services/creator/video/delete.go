package video

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/utils"
)

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
	var fileURLs []string
	// Wrapped in transaction, so we can roll back if it fails, however
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
