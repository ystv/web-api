package encoder

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/encode"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrNoPreset           = errors.New("no presetID on video")
	ErrNoVideoFiles       = errors.New("no video files")
	ErrNoSourceFile       = errors.New("no source file set")
	ErrTooManySourceFiles = errors.New("too many source files set")
	ErrNoFormats          = errors.New("preset has no formats set")
	ErrNoArgs             = errors.New("no encoding arguments set")
)

type Encoder struct {
	encode creator.EncodeRepo
	db     *sqlx.DB
	cdn    *s3.S3
	mq     *amqp.Connection
	c      *http.Client
	conf   *Config
}

type Config struct {
	VTEndpoint  string
	ServeBucket string
}

func NewEncoder(db *sqlx.DB, cdn *s3.S3, mq *amqp.Connection, conf *Config) *Encoder {
	return &Encoder{
		encode: encode.NewStore(db),
		db:     db,
		cdn:    cdn,
		mq:     mq,
		c:      &http.Client{},
		conf:   conf,
	}
}

type (
	VideoItem struct {
		VideoID  int      `db:"video_id"`
		PresetID null.Int `db:"preset_id"`
		Files    []VideoFile
	}
	VideoFile struct {
		FileID         int      `db:"file_id"`
		URI            string   `db:"uri"`
		EncodeFormatID null.Int `db:"encode_format"`
		IsSource       bool     `db:"is_source"`
	}
	EncodeFormat struct {
		Arguments  string `db:"arguments"`
		FileSuffix string `db:"file_suffix"`
	}
)
