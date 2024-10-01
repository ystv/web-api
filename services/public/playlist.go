package public

import (
	"context"
	"fmt"
	"time"
)

// Playlist is a list of videos
// Separate from series and can contain videos from anywhere
type Playlist struct {
	PlaylistID  int         `db:"playlist_id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description string      `db:"description" json:"description"`
	Thumbnail   string      `db:"thumbnail" json:"thumbnail"`
	Videos      []VideoItem `json:"videos"`
}

var _ PlaylistRepo = &Store{}

// GetPlaylist returns a playlist object containing a list of videos and metadata
func (s *Store) GetPlaylist(ctx context.Context, playlistID int) (Playlist, error) {
	var p Playlist

	// Retrieve playlist metadata information
	err := s.db.GetContext(ctx, &p, `
		SELECT playlist_id, name, description, thumbnail
		FROM video.playlists
		WHERE playlist_id = $1;`, playlistID)
	if err != nil {
		return p, fmt.Errorf("failed to get playlist meta: %w", err)
	}

	// Retrieve videos of playlist
	err = s.db.SelectContext(ctx, &p.Videos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.playlist_items vid_list
		INNER JOIN video.items item ON vid_list.video_item_id = item.video_id
		WHERE playlist_id = $1
		ORDER BY position;`, playlistID)
	if err != nil {
		return p, fmt.Errorf("failed to get associated videos: %w", err)
	}

	return p, nil
}

// GetPlaylistPopular returns a playlist of the most popular videos
func (s *Store) GetPlaylistPopular(ctx context.Context, fromPeriod time.Time) (Playlist, error) {
	p := Playlist{
		PlaylistID:  0,
		Name:        "Popular",
		Description: "Popular videos",
	}

	err := s.db.SelectContext(ctx, &p.Videos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.items
		WHERE broadcast_date > $1
		ORDER BY views DESC
		LIMIT 30;`, fromPeriod)
	if err != nil {
		return p, fmt.Errorf("failed to get playlist videos: %w", err)
	}

	return p, nil
}

// GetPlaylistPopularByAllTime returns a playlist of the most popular videos from all time
func (s *Store) GetPlaylistPopularByAllTime(ctx context.Context) (Playlist, error) {
	p := Playlist{
		PlaylistID:  0,
		Name:        "Popular",
		Description: "Popular videos",
	}

	err := s.db.SelectContext(ctx, &p.Videos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.items
		ORDER BY views DESC
		LIMIT 30;`)
	if err != nil {
		return p, fmt.Errorf("failed to get playlist videos: %w", err)
	}

	return p, nil
}

// GetPlaylistPopularByPastYear returns a playlist of the most popular videos from the past year
func (s *Store) GetPlaylistPopularByPastYear(ctx context.Context) (Playlist, error) {
	p := Playlist{
		PlaylistID:  0,
		Name:        "Popular",
		Description: "Popular videos",
	}

	err := s.db.SelectContext(ctx, &p.Videos, `
		SELECT DISTINCT item.video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.items item
		INNER JOIN video.hits hit ON item.video_id = hit.video_id
		WHERE start_time > now() - interval '1 year'
		ORDER BY views DESC
		LIMIT 30;`)
	if err != nil {
		return p, fmt.Errorf("failed to get playlist videos: %w", err)
	}

	return p, nil
}

// GetPlaylistPopularByPastMonth returns a playlist of the most popular videos from the past month
func (s *Store) GetPlaylistPopularByPastMonth(ctx context.Context) (Playlist, error) {
	p := Playlist{
		PlaylistID:  0,
		Name:        "Popular",
		Description: "Popular videos",
	}

	err := s.db.SelectContext(ctx, &p.Videos, `
		SELECT DISTINCT item.video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.items item
		INNER JOIN video.hits hit ON item.video_id = hit.video_id
		WHERE start_time > now() - interval '1 month'
		ORDER BY views DESC
		LIMIT 30;`)
	if err != nil {
		return p, fmt.Errorf("failed to get playlist videos: %w", err)
	}

	return p, nil
}

// GetPlaylistRandom retrieves a random list of videos
func (s *Store) GetPlaylistRandom(ctx context.Context) (Playlist, error) {
	p := Playlist{
		PlaylistID:  0,
		Name:        "Random",
		Description: "Random videos",
	}

	err := s.db.SelectContext(ctx, &p.Videos, `
		SELECT video_id, series_id, name, url, description, thumbnail,
		broadcast_date, views, duration
		FROM video.items TABLESAMPLE system_rows(30);`)

	if err != nil {
		return p, fmt.Errorf("failed to get playlist videos: %w", err)
	}

	return p, nil
}
