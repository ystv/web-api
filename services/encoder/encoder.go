package encoder

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/encode"
)

var (
	ErrNoPreset               = errors.New("no presetID on video")
	ErrNoVideoFiles           = errors.New("no video files")
	ErrNoSourceFile           = errors.New("no source file set")
	ErrTooManySourceFiles     = errors.New("too many source files set")
	ErrNoFormats              = errors.New("preset has no formats set")
	ErrNoArgs                 = errors.New("no encoding arguments set")
	ErrVTFailedToCreate       = errors.New("vt failed to create encode job")
	ErrVTFailedToAuthenticate = errors.New("failed to authenticate to vt")
	ErrVTUnknownResponse      = errors.New("unknown vt response")
	_                         = ErrVTFailedToCreate
)

type Encoder struct {
	encode creator.EncodeRepo
	db     *sqlx.DB
	cdn    *s3.Client
	conf   *Config
}

type Config struct {
	VTEndpoint  string
	ServeBucket string
}

func NewEncoder(db *sqlx.DB, cdn *s3.Client, conf *Config) *Encoder {
	return &Encoder{
		encode: encode.NewStore(db),
		db:     db,
		cdn:    cdn,
		conf:   conf,
	}
}

type (
	VideoItem struct {
		VideoID  int  `db:"video_id"`
		PresetID *int `db:"preset_id"`
		Files    []VideoFile
	}
	VideoFile struct {
		FileID         int    `db:"file_id"`
		EncodeFormatID int    `db:"format_id"`
		URI            string `db:"uri"`
		IsSource       bool   `db:"is_source"`
	}
	EncodeFormat struct {
		Arguments  string `db:"arguments"`
		FileSuffix string `db:"file_suffix"`
	}
	// TaskIdentification is for initially informing the user
	// of their job starting and its given ID for later
	// checking
	TaskIdentification struct {
		State  string `json:"state"`
		TaskID string `json:"taskID"`
	}
)
