package playlist

import (
	"database/sql"

	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

type (
	// Playlist represents a playlist object including the metas of the videos
	Playlist struct {
		Meta
		Videos []video.Meta `json:"videos"`
	}
	// Meta represents the metadata of a playlist
	Meta struct {
		ID          int         `db:"id" json:"id"`
		Name        string      `db:"name" json:"name"`
		Description null.String `db:"description" json:"description"`
		Thumbnail   null.String `db:"thumbnail" json:"thumbnail"`
		Status      string      `db:"status"`
	}
)

// All lists all playlists metadata
func All() ([]Playlist, error) {
	p := []Playlist{}
	err := utils.DB.Select(&p,
		`SELECT id, name, description, thumbnail, status
		FROM video.playlists`)
	return p, err
}

// Get returns a playlist and it's video's
func Get(playlistID int) (Playlist, error) {
	p := Playlist{}
	err := utils.DB.Get(&p,
		`SELECT id, name, description, thumbnail, status
		FROM video.playlists
		WHERE id = $1`, playlistID)
	if err != nil {
		return p, err
	}
	err = utils.DB.Select(&p.Videos,
		`SELECT name, url
		FROM video.items
		INNER JOIN video.playlist_items ON id = video_item_id`)
	return p, err
}

// New makes a playlist item
func New(p Playlist) (sql.Result, error) {
	sql, err := utils.DB.Exec(
		`INSERT INTO video.playlists(name, description, thumbnail, status)
		VALUES ($1, $2, $3, $4);`, p.Name, p.Description, p.Thumbnail, p.Status)
	if err != nil {
		return sql, err
	}
	// TODO finish query
	return utils.DB.Exec(`INSERT INTO video.playlist_items(`)
}
