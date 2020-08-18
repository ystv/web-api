package encode

import (
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// Format represents the encode that is applied
// to a file.
type Format struct {
	FormatID    int         `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
	MimeType    string      `json:"mimeType" db:"mime_type"`
	Mode        string      `json:"mode" db:"mode"`
	Width       int         `json:"width" db:"width"`
	Height      int         `json:"height" db:"height"`
	Watermarked bool        `json:"watermarked" db:"watermarked"`
}

// FormatList lists all encode formats
func FormatList() ([]Format, error) {
	e := []Format{}
	err := utils.DB.Select(&e, `SELECT id, name, description, mime_type, mode, width, height, watermarked
							FROM video.encode_formats;`)
	if err != nil {
		return nil, err
	}
	return e, nil
}
