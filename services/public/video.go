package public

import "time"

type (
	// VideoMeta represents basic information about the videoitem used for listing.
	VideoMeta struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		Description   string    `json:"description"`
		BroadcastDate time.Time `json:"broadcastDate"`
		Views         int       `json:"views"`
		Duration      int       `json:"duration"`
	}
	// VideoFile represents each file that a video item has stored.
	VideoFile struct {
		ID     int    `json:"id"`
		URI    string `json:"uri"`
		Preset string `json:"preset"`
	}
	// VideoItem represents the public basic in-depth information including videofiles.
	VideoItem struct {
		VideoMeta VideoMeta   `json:"videoMeta"`
		Files     []VideoFile `json:"files"`
	}
)

// VideoList returns all video metadata
func VideoList(index int, offset int) ([]VideoMeta, error) {
	videos := []VideoMeta{
		{
			ID: 1, Name: "Vid 1",
		},
		{
			ID: 2, Name: "Vid 2",
		},
	}
	return videos, nil
}

// VideoFind returns a VideoItem, including the files, based on a given VideoItem ID.
func VideoFind(id int) (VideoItem, error) {
	return VideoItem{
		VideoMeta: VideoMeta{
			ID:   1,
			Name: "Video we found",
		},
		Files: []VideoFile{
			{
				ID:     1,
				URI:    "cdn.ystv.co.uk",
				Preset: "480p",
			},
			{
				ID:     2,
				URI:    "cdn.ystv.co.uk",
				Preset: "720p",
			},
		},
	}, nil
}
