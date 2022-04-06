package creator

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/breadcrumb"
	"github.com/ystv/web-api/services/creator/encode"
	"github.com/ystv/web-api/services/creator/playlist"
	"github.com/ystv/web-api/services/creator/playout"
	"github.com/ystv/web-api/services/creator/series"
	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/services/encoder"
	"github.com/ystv/web-api/utils"
)

// Repos represents all our data repositories
type Repos struct {
	access     *utils.Accesser
	video      creator.VideoRepo
	series     creator.SeriesRepo
	playlist   creator.PlaylistRepo
	channel    creator.ChannelRepo
	breadcrumb creator.BreadcrumbRepo
	encode     creator.EncodeRepo
	creator    creator.StatRepo
}

type Config struct {
	IngestBucket string
	ServeBucket  string
}

// NewRepos creates our data repositories
func NewRepos(db *pgx.Conn, cdn *s3.S3, enc *encoder.Encoder, access *utils.Accesser, conf *Config) *Repos {
	config := &creator.Config{
		IngestBucket: conf.IngestBucket,
		ServeBucket:  conf.ServeBucket,
	}
	return &Repos{
		access,
		video.NewStore(db, cdn, enc, config),
		series.NewController(db, cdn, enc, config),
		playlist.NewStore(db),
		playout.NewStore(db),
		breadcrumb.NewController(db, cdn, enc, config),
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
