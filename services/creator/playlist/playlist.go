package playlist

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/utils"
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
		Description null.String `db:"description" json:"description"`
		Thumbnail   null.String `db:"thumbnail" json:"thumbnail"`
		Status      string      `db:"status" json:"status"`
		CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
		CreatedBy   int         `db:"created_by" json:"createdBy"`
	}
)

// All lists all playlists metadata
func All() ([]Playlist, error) {
	p := []Playlist{}
	err := utils.DB.Select(&p,
		`SELECT playlist_id, name, description, thumbnail, status, created_at, created_by
		FROM video.playlists;`)
	return p, err
}

// Get returns a playlist and it's video's
func Get(playlistID int) (Playlist, error) {
	p := Playlist{}
	err := utils.DB.Get(&p,
		`SELECT playlist_id, name, description, thumbnail, status, created_at, created_by
		FROM video.playlists
		WHERE playlist_id = $1;`, playlistID)
	if err != nil {
		return p, err
	}
	err = utils.DB.Select(&p.Videos,
		`SELECT video_id, series_id, name video_name, url, EXTRACT(EPOCH FROM duration)::int AS duration, views, tags, broadcast_date, created_at
		FROM video.items
		INNER JOIN video.playlist_items ON video_id = video_item_id
		ORDER BY position ASC;`)
	return p, err
}

// New makes a playlist item
func New(p Playlist) (sql.Result, error) {
	res, err := utils.DB.Exec(
		`INSERT INTO video.playlists(name, description, thumbnail, status, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6);`, p.Name, p.Description, p.Thumbnail, p.Status, p.CreatedAt, p.CreatedBy)
	if err != nil {
		return res, err
	}
	if len(p.Videos) != 0 {
		res, err = AddVideos(p)
	}
	return res, err
}

// AddVideo adds a single video to a playlist
func AddVideo(p Meta, v *video.Meta) (sql.Result, error) {
	return utils.DB.Exec(`INSERT INTO video.playlist_items (playlist_id, video_item_id) VALUES ($1, $2);`, p.ID, v.ID)
}

// DeleteVideo deletes a single video from a playlist
func DeleteVideo(p Meta, v video.Meta) (sql.Result, error) {
	return utils.DB.Exec(`DELETE FROM video.playlist_items WHERE playlist_id = $1 AND video_item_id = $2;`, p.ID, v.ID)
}

// AddVideos adds multiple videos to a playlist
func AddVideos(p Playlist) (sql.Result, error) {
	var res sql.Result
	txn, err := utils.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := txn.Prepare(pq.CopyIn("video.playlist_items", "playlist_id", "video_item_id"))
	if err != nil {
		return nil, err
	}
	for _, video := range p.Videos {
		res, err = stmt.Exec(p.ID, video.ID)
		if err != nil {
			return res, err
		}
	}
	res, err = stmt.Exec()
	err = stmt.Close()
	if err != nil {
		return res, err
	}
	err = txn.Commit()
	if err != nil {
		return res, err
	}
	return res, err
}
