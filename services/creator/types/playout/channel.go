package playout

import "time"

// This is currently quite bare bones; it is hoped it will integrate with
// ystv/playout to provide more data

// Channel represents a derivative of ystv/playout's channels.
// These are events only rather than linear or events.
type Channel struct {
	URLName        string    `db:"url_name" json:"urlName"`        // "tennis"
	Name           string    `db:"name" json:"name"`               // "YUSU Tennis 2020"
	Description    string    `db:"description" json:"description"` // "Very good tennis"
	Thumbnail      string    `db:"thumbnail" json:"thumbnail"`
	OutputType     string    `db:"output_type" json:"outputType"`
	OutputURL      string    `db:"output_url" json:"outputURL"`
	Visiblity      string    `db:"visibility" json:"visibility"`
	Status         string    `db:"status" json:"status"`     // "live" or "scheduled" or "cancelled" or "finished"
	Location       string    `db:"location" json:"location"` // "Central Hall"
	ScheduledStart time.Time `db:"scheduled_start" json:"scheduledStart"`
	ScheduledEnd   time.Time `db:"scheduled_end" json:"scheduledEnd"`
}
