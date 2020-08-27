package playlist

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/playlist"
	"github.com/ystv/web-api/utils"
)

// Here for validation to ensure we are meeting the interface
var _ creator.PlaylistRepo = &Store{}

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// All lists all playlists metadata
func (m *Store) All(ctx context.Context) ([]playlist.Playlist, error) {
	p := []playlist.Playlist{}
	err := m.db.SelectContext(ctx, &p,
		`SELECT playlist_id, name, description, thumbnail, status, created_at, created_by
		FROM video.playlists;`)
	return p, err
}

// Get returns a playlist and it's video's
func (m *Store) Get(ctx context.Context, playlistID int) (playlist.Playlist, error) {
	p := playlist.Playlist{}
	err := m.db.GetContext(ctx, &p,
		`SELECT playlist_id, name, description, thumbnail, status,
		created_at, created_by, updated_at, updated_by
		FROM video.playlists
		WHERE playlist_id = $1;`, playlistID)
	if err != nil {
		err = fmt.Errorf("failed to select playlist meta: %w", err)
		return p, err
	}
	err = m.db.SelectContext(ctx, &p.Videos,
		`SELECT video_id, series_id, name video_name, url, EXTRACT(EPOCH FROM duration)::int AS duration, views, tags, broadcast_date, created_at
		FROM video.items
		INNER JOIN video.playlist_items ON video_id = video_item_id
		ORDER BY position ASC;`)
	if err != nil {
		err = fmt.Errorf("failed to selected videos: %w", err)
	}
	return p, err
}

// New makes a playlist item
func (m *Store) New(ctx context.Context, p playlist.Playlist) (int, error) {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO video.playlists(name, description, thumbnail, status, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6);`, p.Name, p.Description, p.Thumbnail, p.Status, p.CreatedAt, p.CreatedBy)
	// TODO do we want to use the time here, or what is passed?
	if err != nil {
		err = fmt.Errorf("failed to insert playlist: %w", err)
		return 0, err // Null video ID?
	}
	if len(p.Videos) == 0 {
		return 0, nil
	}
	err = m.AddVideos(ctx, p)
	if err != nil {
		err = fmt.Errorf("failed to add videos to playlist: %w", err)
		return 0, err
	}
	return 0, nil // TODO return playlist ID
}

// AddVideo adds a single video to a playlist
func (m *Store) AddVideo(ctx context.Context, playlistID, videoID int) error {
	_, err := m.db.ExecContext(ctx, `INSERT INTO video.playlist_items (playlist_id, video_item_id) VALUES ($1, $2);`, playlistID, videoID)
	return err
}

// DeleteVideo deletes a single video from a playlist
func (m *Store) DeleteVideo(ctx context.Context, playlistID, videoID int) error {
	_, err := m.db.ExecContext(ctx, `DELETE FROM video.playlist_items WHERE playlist_id = $1 AND video_item_id = $2;`, playlistID, videoID)
	return err
}

// AddVideos adds multiple videos to a playlist
func (m *Store) AddVideos(ctx context.Context, p playlist.Playlist) error {
	// TODO replace this function with the utils transaction wrapper
	txn, err := m.db.Begin()
	if err != nil {
		err = fmt.Errorf("failed to prepare db transaction: %w", err)
		return err
	}
	stmt, err := txn.PrepareContext(ctx, pq.CopyIn("video.playlist_items", "playlist_id", "video_item_id"))
	if err != nil {
		err = fmt.Errorf("failed to prepare statement: %w", err)
		return err
	}
	for _, video := range p.Videos {
		_, err = stmt.ExecContext(ctx, p.ID, video.ID)
		if err != nil {
			err = fmt.Errorf("failed to insert link between playlist and video: %w", err)
			return err
		}
	}
	_, err = stmt.ExecContext(ctx)
	err = stmt.Close()
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Update will update a playlist
// Accepts playlist metadata, video ID's that will be part of the playlist
func (m *Store) Update(ctx context.Context, p playlist.Meta, videoIDs []int) error {
	return utils.Transact(m.db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx,
			`UPDATE video.playlists SET name = $1, description = $2,
			thumbnail = $3, status = $4, updated_at = $5, updated_by = $6
			WHERE playlist_id = $7`, p.Name, p.Description, p.Thumbnail,
			p.Status, time.Now(), p.UpdatedBy)
		if err != nil {
			return fmt.Errorf("failed to update playlist meta: %w", err)
		}
		// Delete old associated videos
		_, err = tx.ExecContext(ctx, `DELETE FROM video.playlist_items
									WHERE playlist_id = $1;`)
		if err != nil {
			return fmt.Errorf("failed to delete old video links: %w", err)
		}
		// No attached videos
		if len(videoIDs) == 0 {
			return nil
		}
		// Insert new video links
		stmt, err := tx.PrepareContext(ctx,
			`INSERT INTO video.playlist_items(playlist_id, video_item_id, position)
		VALUES ($1, $2, $3);`)
		// TODO do we need position? We can get an order sort of for the order
		// of how they where inserted?
		if err != nil {
			return fmt.Errorf("failed to prepare statement to insert videos: %w", err)
		}
		for idx, videoID := range videoIDs {
			_, err := stmt.ExecContext(ctx, p.ID, videoID, idx)
			if err != nil {
				return fmt.Errorf("failed to insert link between playlist and video: %w", err)
			}
		}
		return nil
	})
}
