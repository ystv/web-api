package creator

import (
	"context"
	"fmt"

	"github.com/ystv/web-api/services/creator/types/stats"
)

// GlobalVideo returns general information about a video library
func (m *Store) GlobalVideo(ctx context.Context) (stats.VideoGlobalStats, error) {
	var s stats.VideoGlobalStats

	err := m.db.GetContext(ctx, &s.TotalVideos,
		`SELECT COUNT(*)
		FROM video.items;`)
	if err != nil {
		err = fmt.Errorf("failed to get number of videos: %w", err)
		return s, err
	}

	err = m.db.GetContext(ctx, &s.TotalPendingVideos,
		`SELECT COUNT(*)
		FROM video.items
		WHERE status = 'pending';`)
	if err != nil {
		err = fmt.Errorf("failed to get number of pending videos: %w", err)
		return s, err
	}

	err = m.db.GetContext(ctx, &s.TotalVideoHits,
		`SELECT COUNT(*)
		FROM public.video_hits;`)
	if err != nil {
		err = fmt.Errorf("failed to get number of video hits: %w", err)
		return s, err
	}

	err = m.db.GetContext(ctx, &s.TotalStorageUsed,
		`SELECT SUM(size)
		FROM video.files;`)
	if err != nil {
		err = fmt.Errorf("failed to get size of video files; %w", err)
		return s, err
	}

	return s, nil
}
