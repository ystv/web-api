package video

import (
	"time"

	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"
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
		URI          string   `db:"uri" json:"uri"`
		EncodeFormat string   `db:"name" json:"encodeFormat"`
		Status       string   `db:"status" json:"status"`
		Size         null.Int `db:"size" json:"size"`
		MimeType     string   `db:"mime_type" json:"mimeType"`
	}
	// TODO make null's pointers, so we can omitempty them during JSON marshal

	// Meta represents just the metadata of a video, used for listing.
	Meta struct {
		ID            int            `db:"video_id" json:"id"`
		SeriesID      int            `db:"series_id" json:"seriesID"`
		Name          string         `db:"video_name" json:"name"`
		URL           string         `db:"url" json:"url"`
		Description   string         `db:"description" json:"description,omitempty"` // when listing description isn't included
		Thumbnail     string         `db:"thumbnail" json:"thumbnail"`
		Duration      string         `db:"duration" json:"duration"`
		Views         int            `db:"views" json:"views"`
		Tags          pq.StringArray `db:"tags" json:"tags"`
		Status        string         `db:"status" json:"status"`
		Preset        `json:"preset"`
		BroadcastDate time.Time `db:"broadcast_date" json:"broadcastDate"`
		CreatedAt     time.Time `db:"created_at" json:"createdAt"`
		CreatedBy     User      `json:"createdBy"`
		UpdatedAt     null.Time `db:"updated_at" json:"updatedAt"`
		UpdatedBy     *User     `json:"updatedBy"`
		DeletedAt     null.Time `db:"deleted_at" json:"deletedAt"`
		DeletedBy     *User     `json:"deletedBy"`
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
		PresetID   null.Int    `db:"preset_id" json:"presetID"`
		PresetName null.String `db:"preset_name" json:"name"`
	}
	// New is the basic information to create a video
	New struct {
		FileID        string    `json:"fileID"`
		SeriesID      int       `json:"seriesID" db:"series_id"`
		Name          string    `json:"name" db:"name"`
		URLName       string    `json:"urlName" db:"url"`
		Description   string    `json:"description" db:"description"`
		Tags          []string  `json:"tags" db:"tags"`
		Preset        int       `json:"preset" db:"preset"`
		PublishType   string    `json:"publishType" db:"status"`
		CreatedAt     time.Time `json:"createdAt" db:"created_by"`
		CreatedBy     int       `json:"createdBy" db:"created_by"`
		BroadcastDate time.Time `json:"broadcastDate" db:"broadcast_date"`
	}
)
