package creator

import (
	"time"

	"github.com/ystv/web-api/services/creator/video"
	"gopkg.in/guregu/null.v4"
)

type (
	// PlaylistMeta represents meta information about a playlist, to be used when listing playlists.
	Meta struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		CreationDate time.Time `json:"creationDate"`
	}
	// Playlist represents a playlist object containing the videos.
	Playlist struct {
		Meta
		Videos []video.Item `json:"videos"`
	}
)

// FindPlaylist returns a playlist with nested videoitems and videofiles by ID
func FindPlaylist(ID int) (Playlist, error) {
	return Playlist{
		Meta: Meta{
			ID:           1,
			Name:         "Fun videos",
			CreationDate: time.Now(),
		},
		Videos: []video.Item{
			{
				ID:       1,
				Name:     "Dog barks",
				Duration: null.IntFrom(400),
				Files: []video.File{
					{
						URI:          "cdn.ystv.co.uk",
						EncodeFormat: "Original",
					}, {
						URI:          "cdn.ystv.co.uk",
						EncodeFormat: "240p"},
				}},
			{ID: 2, Name: "Cat meows"},
			{ID: 3, Name: "Cow moo"},
		},
	}, nil
}

// ListPlaylist returns all playlists
func ListPlaylist() ([]Meta, error) {
	return []Meta{
		{ID: 1, Name: "Fun videos", CreationDate: time.Now()},
		{ID: 2, Name: "Sad videos", CreationDate: time.Now()},
	}, nil
}
