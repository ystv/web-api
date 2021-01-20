package video

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/creator"
)

// TODO update schema so duration is not null

var _ creator.VideoRepo = &Store{}

// Store encapsulates our dependencies
type Store struct {
	db   *sqlx.DB
	cdn  *s3.S3
	conf *creator.Config
}

func getSeason(t time.Time) string {
	m := int(t.Month())
	switch {
	case m >= 9 && m <= 12:
		return "aut"
	case m >= 1 && m <= 6:
		return "spr"
	default:
		return "sum"
	}
}
