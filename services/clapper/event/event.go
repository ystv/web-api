package event

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/clapper"
	"github.com/ystv/web-api/utils"
)

// Store encapsulates our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates our data store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db}
}

// Here to verify we are meeting the interface
var _ clapper.EventRepo = &Store{}

// ListMonth Lists all event meta's for a month
func (m *Store) ListMonth(ctx context.Context, year, month int) (*[]clapper.Event, error) {
	e := []clapper.Event{}
	err := m.db.SelectContext(ctx, &e,
		`SELECT event_id, event_type, name, start_date, end_date, description,
		location, is_private, is_cancelled, is_tentative
		FROM event.events
		WHERE EXTRACT(YEAR FROM start_date) = $1 AND
		EXTRACT(MONTH FROM start_date) = $2;`, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to list month: %w", err)
	}
	return &e, nil
}

// Get returns an event including the signup sheets
func (m *Store) Get(ctx context.Context, eventID int) (*clapper.Event, error) {
	e := clapper.Event{}
	err := m.db.GetContext(ctx, &e,
		`SELECT event_id, event_type, name, start_date, end_date, description,
		location, is_private, is_cancelled, is_tentative
		FROM event.events
		WHERE event_id = $1;`, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event meta: %w", err)
	}
	if e.EventType != "show" {
		err := m.db.SelectContext(ctx, &e.Attendees,
			`SELECT users.user_id, nickname, first_name, last_name, attend_status
			FROM event.attendees attendees
			INNER JOIN people.users users ON attendees.user_id = users.user_id
			WHERE event_id = $1;`, e.EventID)
		if err != nil {
			return nil, fmt.Errorf("failed to get attendees: %w", err)
		}
		return &e, nil
	}
	err = m.db.SelectContext(ctx, &e.Signups,
		`SELECT signup_id, title, description, unlock_date, arrival_time, start_time, end_time
		FROM event.signups
		WHERE event_id = $1;`, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signup sheets: %w", err)
	}
	for i := range e.Signups {
		err := m.db.SelectContext(ctx, &e.Signups[i].Crew,
			`SELECT crew_id, crew.user_id, nickname, first_name, last_name, locked,
			event.positions.position_id, name, description, admin, credited, permission_id
			FROM event.crews crew
			INNER JOIN event.positions ON event.positions.position_id = crew.position_id
			INNER JOIN people.users ON crew.user_id  = people.users.user_id
			WHERE signup_id = $1
			ORDER BY ordering;`, e.Signups[i].SignupID)
		if err != nil {
			return nil, fmt.Errorf("failed to get crew for signup sheets: %w", err)
		}
	}
	return &e, nil
}

// New creates a new event returning the event ID
func (m *Store) New(ctx context.Context, e *clapper.NewEvent, userID int) (int, error) {
	eventID := 0
	err := m.db.QueryRowContext(ctx, `INSERT INTO event.events 
	(event_type, name, start_date, end_date, description, location,
	is_private, is_tentative, created_at, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING event_id;`,
		&e.EventType, &e.Name, &e.StartDate, &e.EndDate, &e.Description,
		&e.Location, &e.IsPrivate, &e.IsTentative, time.Now(), userID).Scan(&eventID)
	if err != nil {
		return eventID, fmt.Errorf("failed to insert new event: %w", err)
	}
	return eventID, nil
}

// Update an existing event
func (m *Store) Update(ctx context.Context, e *clapper.Event, userID int) error {
	eventType := ""
	err := utils.Transact(m.db, func(tx *sqlx.Tx) error {
		err := tx.QueryRowContext(ctx, `UPDATE event.events
		SET event_type = $1, name = $2, start_date = $3, end_date = $4,
		description = $5, location = $6, is_private = $7,
		is_tentative = $8, updated_at = $9, updated_by = $10
		WHERE event_id = $11
		RETURNING
		(SELECT event_id
		FROM event.events
		WHERE event_id = $11);`, &e.EventType, &e.Name, &e.StartDate,
			&e.EndDate, &e.Description, &e.Location, &e.IsPrivate,
			&e.IsTentative, time.Now(), userID, &e.EventID).Scan(&eventType)
		if err != nil {
			return fmt.Errorf("failed to update event meta")
		}
		// Check if the event type is changed
		if eventType == e.EventType {
			return nil
		}
		// We've had a change
		switch e.EventType {
		case "show":
			// other -> show
			_, err = tx.ExecContext(ctx, `DELETE FROM event.signups WHERE event_id = $1;`, e.EventID)
			if err != nil {
				return fmt.Errorf("failed to delete event signups: %w", err)
			}
		default:
			// show -> other
			_, err = tx.ExecContext(ctx, `DELETE FROM event.attendees WHERE event_id = $1;`, e.EventID)
			if err != nil {
				return fmt.Errorf("failed to delete event attendees: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

// Delete an event by EventID
func (m *Store) Delete(ctx context.Context, eventID int) error {
	// The cascade should delete either the attendees or signups
	_, err := m.db.ExecContext(ctx, `DELETE FROM event.events
							WHERE event_id = $1;`, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}
