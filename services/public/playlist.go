package public

import (
	"context"
	"fmt"
)

// Playlist is a list of videos
// Seperate from series and can contain videos from anywhere
type Playlist struct {
	PlaylistID  int         `db:"playlist_id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description string      `db:"description" json:"description"`
	Thumbnail   string      `db:"thumbnail" json:"thumbnail"`
	Videos      []VideoItem `json:"videos"`
}

var _ PlaylistRepo = &Store{}

// GetPlaylist returns a playlist object containing a list of videos and metadata
func (m *Store) GetPlaylist(ctx context.Context, playlistID int) (Playlist, error) {
	p := Playlist{}
	err := m.db.GetContext(ctx, &p, `
		SELECT playlist_id, name, description, thumbnail
		FROM video.playlists
		WHERE playlist_id = $1;`, playlistID)
	if err != nil {
		return p, fmt.Errorf("failed to get playlist meta: %w", err)
	}
	err = m.db.SelectContext(ctx, &p.Videos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.playlist_items vid_list
		INNER JOIN video.items item ON vid_list.video_item_id = item.video_id
		WHERE playlist_id = $1
		ORDER BY position;`, playlistID)
	if err != nil {
		return p, fmt.Errorf("failed to get associated videos: %w", err)
	}
	return p, nil
}
