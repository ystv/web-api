package creator

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/breadcrumb"
	"github.com/ystv/web-api/services/creator/encode"
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
	encode     creator.EncodeRepo
	creator    creator.StatRepo
}

// NewRepos creates our data repositories
func NewRepos(db *sqlx.DB, cdn *s3.S3) *Repos {
	config := &creator.Config{
		IngestBucket: os.Getenv("WAPI_BUKCET_VOD_INGEST"),
		ServeBucket:  os.Getenv("WAPI_BUCKET_VOD_SERVE"),
	}
	return &Repos{
		video.NewStore(db, cdn, config),
		series.NewController(db, cdn, config),
		playlist.NewStore(db),
		breadcrumb.NewController(db, cdn, config),
		encode.NewStore(db),
		creator.NewStore(db),
	}
}

// Stats handles sending general stats about the video library
// @Summary Get global video library information
// @Description Gets the statistics about the global video library.
// @ID get-creator-glob-stats
// @Tags creator
// @Produce json
// @Success 200 {object} stats.VideoGlobalStats
// @Router /v1/internal/creator/stats [get]
func (r *Repos) Stats(c echo.Context) error {
	s, err := r.creator.GlobalVideo(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("stats failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, s)
}
