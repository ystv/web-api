package creator

import "time"

type PlaylistMeta struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

type Playlist struct {
	Meta   PlaylistMeta
	videos []int
}
