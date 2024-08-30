package encoder

import (
	"context"
	"fmt"
)

func (e *Encoder) TranscodeFinished(ctx context.Context, taskID string) error {
	fileID := 0
	err := e.db.GetContext(ctx, &fileID, `
		SELECT file_id
		FROM video.files
		WHERE status = $1;`, fmt.Sprintf("processing/%s", taskID))
	if err != nil {
		return fmt.Errorf("failed to get video files: %w", err)
	}

	_, err = e.db.ExecContext(ctx, `UPDATE video.files SET status = 'public' WHERE file_id = $1;`, fileID)
	if err != nil {
		return fmt.Errorf("failed to update video file")
	}

	return nil
}
