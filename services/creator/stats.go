package creator

import (
	"context"

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
		return nil, err
	}
	err = utils.DB.GetContext(ctx, &s.TotalPendingVideos,
		`SELECT COUNT(*)
		FROM video.items
		WHERE status = 'pending';`)
	err = utils.DB.GetContext(ctx, &s.TotalVideoHits,
		`SELECT COUNT(*)
		FROM public.video_hits;`)
	if err != nil {
		return nil, err
	}
	err = utils.DB.GetContext(ctx, &s.TotalStorageUsed,
		`SELECT SUM(size)
		FROM video.files;`)
	if err != nil {
		return nil, err
	}
	return s, nil
}
