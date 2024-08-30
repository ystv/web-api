package video

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type (
	// Item represents a more readable VideoItem with
	// an array of associated VideoFiles.
	Item struct {
		Meta
		Files []File `db:"files" json:"files"`
	}
	// File represents a more readable VideoFile.
	File struct {
		URI          string `db:"uri" json:"uri"`
		EncodeFormat string `db:"name" json:"encodeFormat"`
		Status       string `db:"status" json:"status"`
		Size         *int   `db:"size" json:"size"`
		MimeType     string `db:"mime_type" json:"mimeType"`
	}

	// Meta represents just the metadata of a video, used for listing.
	Meta struct {
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
		Preset        `json:"preset"`
		BroadcastDate time.Time  `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
		CreatedByID   int        `db:"created_by_id" json:"createdByID"`
		CreatedByNick string     `db:"created_by_nick" json:"createdByNick"`
		UpdatedAt     *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedByID   *int       `db:"updated_by_nick" json:"updatedByID,omitempty"`
		UpdatedByNick *string    `db:"updated_by_nick" json:"updatedByNick,omitempty"`
		DeletedAt     *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedByID   *int       `db:"deleted_by_id" json:"deleteByID,omitempty"`
		DeletedByNick *string    `db:"deleted_by_nick" json:"deleteByNick,omitempty"`
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
	// Preset represents the name and ID of a preset
	Preset struct {
		PresetID   *int    `db:"preset_id" json:"presetID"`
		PresetName *string `db:"preset_name" json:"name"`
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

func (t Tag) Value() (driver.Value, error) {
	if len(t) == 0 {
		return "{}", nil
	}
	return fmt.Sprintf(`{"%s"}`, strings.Join(t, `","`)), nil
}

func (t *Tag) Scan(src interface{}) (err error) {
	var tags []string
	switch src.(type) {
	case string:
		fmt.Println("string")
		err = json.Unmarshal([]byte(src.(string)), &tags)
	case []byte:
		fmt.Println("[]byte")
		fmt.Println(string(src.([]byte)))
		temp := string(src.([]byte))

		temp = strings.TrimLeft(temp, "{")
		temp = strings.TrimRight(temp, "}")
		tags = strings.Split(temp, ",")
	default:
		return errors.New("incompatible type for Tag")
	}
	if err != nil {
		return
	}
	*t = tags
	return nil
}
