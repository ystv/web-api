package encode

import "gopkg.in/guregu/null.v4"

// Preset represents a group of encode formats. A video has a preset applied to it so
// it can generate the video files so a video item.
type (
	Preset struct {
		PresetID    int         `json:"id" db:"id"`
		Name        string      `json:"name" db:"name"`
		Description null.String `json:"description" db:"description"`
		Formats     []Format    `json:"formats"`
	}
	// Format represents the encode that is applied
	// to a file.
	Format struct {
		FormatID    int         `json:"id" db:"id"`
		Name        string      `json:"name" db:"name"`
		Description null.String `json:"description" db:"description"`
		MimeType    string      `json:"mimeType" db:"mime_type"`
		Mode        string      `json:"mode" db:"mode"`
		Width       int         `json:"width" db:"width"`
		Height      int         `json:"height" db:"height"`
		Watermarked bool        `json:"watermarked" db:"watermarked"`
	}
)
