package creator

import (
	"context"

	"github.com/ystv/web-api/utils"
)

type SQLStats struct {
	TotalVideoHits int `db:"hits" json:"hits"`
}

// Stats returns general information about video library
func Stats(ctx context.Context) (*SQLStats, error) {
	s := &SQLStats{}
	err := utils.DB.GetContext(ctx, &s, `SELECT`)
	if err != nil {
		return nil, err
	}
	return s, nil
}
