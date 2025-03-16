package video

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v4"
)

type (
	// ItemDB represents a more readable VideoItem with
	// an array of associated VideoFiles.
	ItemDB struct {
		MetaDB
		Files []FileDB `db:"files" json:"files"`
	}

	// Item represents a more readable VideoItem with
	// an array of associated VideoFiles.
	Item struct {
		Meta
		Files []File `db:"files" json:"files"`
	}

	// FileDB represents a more readable VideoFile.
	FileDB struct {
		URI          string   `db:"uri"`
		EncodeFormat string   `db:"name"`
		Status       string   `db:"status"`
		Size         null.Int `db:"size"`
		MimeType     string   `db:"mime_type"`
	}

	// File represents a more readable VideoFile.
	File struct {
		URI          string `json:"uri"`
		EncodeFormat string `json:"encodeFormat"`
		Status       string `json:"status"`
		Size         *int64 `json:"size,omitempty"`
		MimeType     string `json:"mimeType"`
	}

	// MetaDB represents just the metadata of a video, used for listing.
	MetaDB struct {
		ID            int    `db:"video_id" json:"id"`
		SeriesID      int    `db:"series_id" json:"seriesID"`
		Name          string `db:"video_name" json:"name"`
		URL           string `db:"url" json:"url"`
		Description   string `db:"description" json:"description,omitempty"` // when listing description isn't included
		Thumbnail     string `db:"thumbnail" json:"thumbnail"`
		Duration      int    `db:"duration" json:"duration"`
		Views         int    `db:"views" json:"views"`
		Tags          Tag    `db:"tags" json:"tags"`
		Status        string `db:"status" json:"status"`
		PresetDB      `json:"preset"`
		BroadcastDate time.Time   `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt     time.Time   `db:"created_at" json:"createdAt"`
		CreatedByID   int         `db:"created_by_id" json:"createdByID"`
		CreatedByNick string      `db:"created_by_nick" json:"createdByNick"`
		UpdatedAt     null.Time   `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedByID   null.Int    `db:"updated_by_id" json:"updatedByID,omitempty"`
		UpdatedByNick null.String `db:"updated_by_nick" json:"updatedByNick,omitempty"`
		DeletedAt     null.Time   `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedByID   null.Int    `db:"deleted_by_id" json:"deleteByID,omitempty"`
		DeletedByNick null.String `db:"deleted_by_nick" json:"deleteByNick,omitempty"`
	}

	// Meta represents just the metadata of a video, used for listing.
	Meta struct {
		ID            int    `json:"id"`
		SeriesID      int    `json:"seriesID"`
		Name          string `json:"name"`
		URL           string `json:"url"`
		Description   string `json:"description,omitempty"` // when listing description isn't included
		Thumbnail     string `json:"thumbnail"`
		Duration      int    `json:"duration"`
		Views         int    `json:"views"`
		Tags          Tag    `json:"tags"`
		Status        string `json:"status"`
		Preset        `json:"preset"`
		BroadcastDate time.Time  `json:"broadcastDate"`
		CreatedAt     time.Time  `json:"createdAt"`
		CreatedByID   int        `json:"createdByID"`
		CreatedByNick string     `json:"createdByNick"`
		UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
		UpdatedByID   *int64     `json:"updatedByID,omitempty"`
		UpdatedByNick *string    `json:"updatedByNick,omitempty"`
		DeletedAt     *time.Time `json:"deletedAt,omitempty"`
		DeletedByID   *int64     `json:"deleteByID,omitempty"`
		DeletedByNick *string    `json:"deleteByNick,omitempty"`
	}

	// MetaCal represents simple metadata for a calendar
	MetaCal struct {
		ID            int    `db:"video_id" json:"id"`
		Name          string `db:"name" json:"name"`
		Status        string `db:"status" json:"status"`
		BroadcastDate string `db:"broadcast_date" json:"broadcastDate"`
	}

	// User represents the nickname and ID of a user
	User struct {
		UserID   int    `db:"user_id" json:"userID"`
		Nickname string `db:"nickname" json:"userNickname"`
	}

	// PresetDB represents the name and ID of a preset
	PresetDB struct {
		PresetID   null.Int    `db:"preset_id" json:"presetID"`
		PresetName null.String `db:"preset_name" json:"name"`
	}

	// Preset represents the name and ID of a preset
	Preset struct {
		PresetID   *int64  `json:"presetID,omitempty"`
		PresetName *string `json:"name,omitempty"`
	}

	// New is the basic information to create a video
	New struct {
		FileID        string    `json:"fileID"`
		SeriesID      int       `json:"seriesID" db:"series_id"`
		Name          string    `json:"name" db:"name"`
		URLName       string    `json:"urlName" db:"url"`
		Description   string    `json:"description" db:"description"`
		Tags          []string  `json:"tags" db:"tags"`
		PresetID      int       `json:"presetID" db:"preset_id"`
		PublishType   string    `json:"publishType" db:"status"`
		CreatedAt     time.Time `json:"createdAt" db:"created_by"`
		CreatedBy     int       `json:"createdBy" db:"created_by"`
		BroadcastDate time.Time `json:"broadcastDate" db:"broadcast_date"`
	}

	Tag []string
)

var (
	ErrNotFound = errors.New("video not found")
)

func (t *Tag) Value() (driver.Value, error) {
	if len(*t) == 0 {
		return "{}", nil
	}
	return fmt.Sprintf(`{"%s"}`, strings.Join(*t, `","`)), nil
}

func (t *Tag) Scan(src interface{}) (err error) {
	var tags []string
	switch s := src.(type) {
	case string:
		err = json.Unmarshal([]byte(s), &tags)
		if err != nil {
			return
		}
	case []byte:
		temp := string(s)
		temp = strings.TrimLeft(temp, "{")
		temp = strings.TrimRight(temp, "}")
		tags = strings.Split(temp, ",")
	default:
		return errors.New("incompatible type for Tag")
	}

	*t = tags
	return nil
}
