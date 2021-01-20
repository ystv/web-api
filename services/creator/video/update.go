package video

import (
	"context"

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
	// Need to check if preset has changed
	presetID := null.Int{}
	s.db.GetContext(ctx, &presetID, `
		UPDATE video.items SET
		series_id = $1,
		name = $2,
		url = $3,
		description = $4
		thumbnail = $5,
		duration = (EPOCH FROM $6),
		genre  = $7,
		tags = $8,
		series_position = $9,
		status = $10,
		preset = $11,
		broadcast_date = $12,
		updated_at = $13,
		updated_by = $14
		
		RETURNING preset;`)
	if m.Preset.PresetID != presetID {
		// preset change, need to schedule new video files
	}
	return nil
}
