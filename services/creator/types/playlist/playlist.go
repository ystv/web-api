package playlist

import (
	"time"

	"github.com/ystv/web-api/services/creator/types/video"
	"gopkg.in/guregu/null.v4"
)

type (
	// Playlist represents a playlist object including the metas of the videos
	Playlist struct {
		Meta
		Videos []video.Meta `json:"videos,omitempty"`
	}
	// Meta represents the metadata of a playlist
	Meta struct {
		ID          int         `db:"playlist_id" json:"id"`
		Name        string      `db:"name" json:"name"`
		Description string      `db:"description" json:"description"`
		Thumbnail   null.String `db:"thumbnail" json:"thumbnail"`
		Status      string      `db:"status" json:"status"`
		CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
		CreatedBy   int         `db:"created_by" json:"createdBy"`
		UpdatedAt   null.Time   `db:"updated_at" json:"updatedAt"`
		UpdatedBy   null.Int    `db:"updated_by" json:"updatedBy"`
	}
)
