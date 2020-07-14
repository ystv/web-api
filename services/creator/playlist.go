package creator

import "time"

type (
	// PlaylistMeta represents meta information about a playlist, to be used when listing playlists.
	PlaylistMeta struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		CreationDate time.Time `json:"creationDate"`
	}
	// Playlist represents a playlist object containing the videos.
	Playlist struct {
		ID           int         `json:"id"`
		Name         string      `json:"name"`
		CreationDate time.Time   `json:"creationDate"`
		Videos       []VideoItem `json:"videos"`
	}
)

// FindPlaylist returns a playlist with nested videoitems and videofiles by ID
func FindPlaylist(ID int) (Playlist, error) {
	return Playlist{ID: 1,
		Name:         "Fun videos",
		CreationDate: time.Now(),
		Videos: []VideoItem{
			{
				ID:       1,
				Name:     "Dog barks",
				Duration: 400,
				Files: []VideoFile{
					{
						ID:     1,
						URI:    "cdn.ystv.co.uk",
						Preset: "Original",
					}, {
						ID:     2,
						URI:    "cdn.ystv.co.uk",
						Preset: "240p"},
				}},
			{ID: 2, Name: "Cat meows"},
			{ID: 3, Name: "Cow moo"},
		}}, nil
}

// ListPlaylist returns all playlists
func ListPlaylist() ([]PlaylistMeta, error) {
	return []PlaylistMeta{
		{ID: 1, Name: "Fun videos", CreationDate: time.Now()},
		{ID: 2, Name: "Sad videos", CreationDate: time.Now()},
	}, nil
}
