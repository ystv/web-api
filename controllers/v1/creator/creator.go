package creator

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
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
type (
	Repos interface {
		Stats(c echo.Context) error
		EncodeRepo
		PlaylistRepo
		PlayoutRepo
		SeriesRepo
		VideoRepo
	}

	EncodeRepo interface {
		ListEncodeFormats(c echo.Context) error
		NewEncodeFormat(c echo.Context) error
		UpdateEncodeFormat(c echo.Context) error
		DeleteEncodeFormat(c echo.Context) error
		ListEncodePresets(c echo.Context) error
		NewEncodePreset(c echo.Context) error
		UpdateEncodePreset(c echo.Context) error
		DeleteEncodePreset(c echo.Context) error
	}

	PlaylistRepo interface {
		ListPlaylists(c echo.Context) error
		GetPlaylist(c echo.Context) error
		NewPlaylist(c echo.Context) error
		UpdatePlaylist(c echo.Context) error
		DeletePlaylist(c echo.Context) error
	}

	PlayoutRepo interface {
		ListChannels(c echo.Context) error
		NewChannel(c echo.Context) error
		UpdateChannel(c echo.Context) error
		DeleteChannel(c echo.Context) error
	}

	SeriesRepo interface {
		ListSeries(c echo.Context) error
		GetSeries(c echo.Context) error
		UpdateSeries(c echo.Context) error
		DeleteSeries(c echo.Context) error
	}

	VideoRepo interface {
		GetVideo(c echo.Context) error
		NewVideo(c echo.Context) error
		UpdateVideoMeta(c echo.Context) error
		DeleteVideo(c echo.Context) error
		ListVideos(c echo.Context) error
		ListVideosByUser(c echo.Context) error
		ListVideosByMonth(c echo.Context) error
		SearchVideo(c echo.Context) error
	}

	Store struct {
		access     utils.Repo
		video      creator.VideoRepo
		series     creator.SeriesRepo
		playlist   creator.PlaylistRepo
		channel    creator.ChannelRepo
		breadcrumb creator.BreadcrumbRepo
		encode     creator.EncodeRepo
		creator    creator.StatRepo
	}

	Config struct {
		IngestBucket string
		ServeBucket  string
	}
)

// NewRepos creates our data repositories
func NewRepos(db *sqlx.DB, cdn *s3.S3, enc encoder.Repo, access utils.Repo, conf *Config, cdnEndpoint string) Repos {
	config := &creator.Config{
		IngestBucket: conf.IngestBucket,
		ServeBucket:  conf.ServeBucket,
		Endpoint:     cdnEndpoint,
	}
	return &Store{
		access,
		video.NewStore(db, cdn, enc, config),
		series.NewController(db, cdn, enc, config),
		playlist.NewStore(db),
		playout.NewStore(db, cdn, config),
		breadcrumb.NewController(db, cdn, enc, config),
		encode.NewStore(db),
		creator.NewStore(db),
	}
}

// Stats handle sending general stats about the video library
// @Summary Get global video library information
// @Description Gets the statistics about the global video library.
// @ID get-creator-glob-stats
// @Tags creator
// @Produce json
// @Success 200 {object} stats.VideoGlobalStats
// @Router /v1/internal/creator/stats [get]
func (s *Store) Stats(c echo.Context) error {
	stats, err := s.creator.GlobalVideoStats(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("stats failed: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, stats)
}
