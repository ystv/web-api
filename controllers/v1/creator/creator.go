package creator

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/breadcrumb"
	"github.com/ystv/web-api/services/creator/playlist"
	"github.com/ystv/web-api/services/creator/series"
	"github.com/ystv/web-api/services/creator/video"
)

// Repos represents all our data repositories
type Repos struct {
	video      creator.VideoRepo
	series     creator.SeriesRepo
	playlist   creator.PlaylistRepo
	breadcrumb creator.BreadcrumbRepo
}

// NewRepos creates our data repositories
func NewRepos(db *sqlx.DB, cdn *s3.S3) *Repos {
	return &Repos{
		video.NewStore(db, cdn),
		series.NewController(db, cdn),
		playlist.NewStore(db),
		breadcrumb.NewController(db, cdn),
	}
}

// Stats handles sending general stats about the video library
func Stats(c echo.Context) error {
	s, err := creator.Stats(context.Background())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, s)
}
