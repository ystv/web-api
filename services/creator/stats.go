package creator

import (
	"context"
	"fmt"

	"github.com/ystv/web-api/utils"
)

type SQLStats struct {
	TotalVideos        int `db:"videos" json:"totalVideos"`
	TotalPendingVideos int `db:"pending" json:"totalPendingVideos"`
	TotalVideoHits     int `db:"hits" json:"totalVideoHits"`
	TotalStorageUsed   int `db:"storage" json:"totalStorageUsed"`
}

// Stats returns general information about video library
func Stats(ctx context.Context) (*SQLStats, error) {
	s := &SQLStats{}
	err := utils.DB.GetContext(ctx, &s.TotalVideos,
		`SELECT COUNT(*)
		FROM video.items;`)
	if err != nil {
		err = fmt.Errorf("failed to get number of videos: %w", err)
		return nil, err
	}
	err = utils.DB.GetContext(ctx, &s.TotalPendingVideos,
		`SELECT COUNT(*)
		FROM video.items
		WHERE status = 'pending';`)
	if err != nil {
		err = fmt.Errorf("failed to get number of pending videos: %w", err)
		return nil, err
	}
	err = utils.DB.GetContext(ctx, &s.TotalVideoHits,
		`SELECT COUNT(*)
		FROM public.video_hits;`)
	if err != nil {
		err = fmt.Errorf("failed to get number of video hits: %w", err)
		return nil, err
	}
	err = utils.DB.GetContext(ctx, &s.TotalStorageUsed,
		`SELECT SUM(size)
		FROM video.files;`)
	if err != nil {
		err = fmt.Errorf("failed to get size of video files; %w", err)
		return nil, err
	}
	return s, nil
}
