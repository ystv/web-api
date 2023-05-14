package public

import (
	"context"
	"fmt"
	"time"
)

var _ StreamRepo = &Store{}

// This is currently quite barebones, it is hoped it will integrate with
// ystv/playout in order to provide more data

// Channel represents a derivative of ystv/playout's channels.
// These are event only rather than linear or event.
type Channel struct {
	URLName        string    `db:"url_name" json:"urlName"`        // "tennis"
	Name           string    `db:"name" json:"name"`               // "YUSU Tennis 2020"
	Description    string    `db:"description" json:"description"` // "Very good tennis"
	Thumbnail      string    `db:"thumbnail" json:"thumbnail"`
	OutputType     string    `db:"output_type" json:"outputType"`
	OutputURL      string    `db:"output_url" json:"outputURL"`
	Status         string    `db:"status" json:"status"`     // "live" or "scheduled" or "cancelled" or "finished"
	Location       string    `db:"location" json:"location"` // "Central Hall"
	ScheduledStart time.Time `db:"scheduled_start" json:"scheduledStart"`
	ScheduledEnd   time.Time `db:"scheduled_end" json:"scheduledEnd"`
}

// ListChannels will list all public channels
func (s *Store) ListChannels(ctx context.Context) ([]Channel, error) {
	var chs []Channel
	err := s.db.SelectContext(ctx, &chs, `
		SELECT url_name, name, description, thumbnail, output_type, output_url,
		status, location, scheduled_start, scheduled_end
		FROM playout.channel
		WHERE visibility = 'public';`)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	return chs, nil
}

// GetChannel will get a public or unlisted channel
func (s *Store) GetChannel(ctx context.Context, urlName string) (Channel, error) {
	ch := Channel{}
	err := s.db.GetContext(ctx, &ch, `
		SELECT url_name, name, description, thumbnail, output_type, output_url,
		status, location, scheduled_start, scheduled_end
		FROM playout.channel
		WHERE visibility IN ('public', 'unlisted')
		AND url_name = $1;`, urlName)
	if err != nil {
		return Channel{}, fmt.Errorf("failed to get channels: %w", err)
	}
	return ch, nil
}
