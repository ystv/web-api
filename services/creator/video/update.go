package video

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ystv/web-api/services/creator/types/video"
	"gopkg.in/guregu/null.v4"
)

// UpdateMeta updates a video's metadata
//
// This wont update:
// * duration
// * views
//
func (s *Store) UpdateMeta(ctx context.Context, m video.Meta) error {
	_, err := s.cdn.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.conf.ServeBucket),
		CopySource: aws.String(s.conf.IngestBucket + "/" + m.Thumbnail),
		Key:        aws.String(m.Thumbnail),
	})
	if err != nil {
		return fmt.Errorf("failed to copy thumbnail: %w", err)
	}

	m.Thumbnail = "https://cdn.ystv.co.uk/" + s.conf.ServeBucket + "/" + m.Thumbnail

	// Need to check if preset has changed
	presetID := null.Int{}
	err = s.db.GetContext(ctx, &presetID, `
		UPDATE video.items SET
		series_id = $1,
		name = $2,
		url = $3,
		description = $4,
		thumbnail = $5,
		tags = $6,
		status = $7,
		preset = $8,
		broadcast_date = $9,
		updated_at = $10,
		updated_by = $11
		
		WHERE video_id = $12

		RETURNING preset;`, m.SeriesID, m.Name, m.URL, m.Description, m.Thumbnail,
		m.Tags, m.Status, m.Preset.PresetID, m.BroadcastDate, m.UpdatedAt, m.UpdatedByID, m.ID)
	if err != nil {
		return fmt.Errorf("failed to update video in db: %w", err)
	}
	if m.Preset.PresetID != presetID {
		// preset change, need to schedule new video files
		s.enc.RefreshVideo(ctx, m.ID)
	}
	return nil
}