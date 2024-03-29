package playout

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/playout"
)

// Here for validation to ensure we are meeting the interface
var _ creator.ChannelRepo = &Store{}

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// ListChannels list all channels
func (s *Store) ListChannels(ctx context.Context) ([]playout.Channel, error) {
	var chs []playout.Channel
	err := s.db.SelectContext(ctx, &chs, `
		SELECT url_name, name, description, thumbnail, output_type, output_url,
		visibility, status, location, scheduled_start, scheduled_end
		FROM playout.channel;`)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	return chs, nil
}

// NewChannel create a new channel
func (s *Store) NewChannel(ctx context.Context, ch playout.Channel) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO playout.channel
		(url_name, name, description, thumbnail, output_type, output_url,
		visibility, status, location, scheduled_start, scheduled_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		ch.URLName, ch.Name, ch.Description, ch.Thumbnail, ch.OutputType,
		ch.OutputURL, ch.Visibility, ch.Status, ch.Location, ch.ScheduledStart,
		ch.ScheduledEnd)
	if err != nil {
		return fmt.Errorf("failed to create channe: %w", err)
	}
	return nil
}

// UpdateChannel update a channel
func (s *Store) UpdateChannel(ctx context.Context, ch playout.Channel) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE playout.channel SET
			url_name = $1, name = $2, description = $3, thumbnail = $4,
			output_type = $5, output_url = $6, visibility = $7,	status = $8,
			location = $9, scheduled_start = $10, scheduled_end = $11
		WHERE url_name = $1;`,
		ch.URLName, ch.Name, ch.Description, ch.Thumbnail, ch.OutputType,
		ch.OutputURL, ch.Visibility, ch.Status, ch.Location, ch.ScheduledStart,
		ch.ScheduledEnd)
	if err != nil {
		return fmt.Errorf("failed to update channel: %w", err)
	}
	return nil
}

// DeleteChannel delete a channel
func (s *Store) DeleteChannel(ctx context.Context, urlName string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM playout.channel
		WHERE url_name = $1;`, urlName)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}
	return nil
}
