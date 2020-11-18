package signup

import (
	"context"
	"fmt"

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
var _ clapper.SignupRepo = &Store{}

// New creates a new signup sheet
func (m *Store) New(ctx context.Context, eventID int, s clapper.Signup) (int, error) {
	signupID := 0
	err := utils.Transact(m.db, func(tx *sqlx.Tx) error {
		err := tx.QueryRowContext(ctx, `INSERT INTO event.signups
		(event_id, title, description, unlock_date, arrival_time, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING signup_id;`,
			eventID, s.Title, s.Description, s.UnlockDate, s.ArrivalTime, s.StartTime, s.EndTime).Scan(&signupID)
		if err != nil {
			return fmt.Errorf("failed to insert new signup shee metat: %w", err)
		}
		// Check if positions have been added
		// TODO I'm not too sure on using the signup struct for this,
		// maybe another input variable instead?
		if len(s.Crew) == 0 {
			return nil
		}
		s.SignupID = signupID
		err = m.addCrew(ctx, tx, s)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert new signup sheet: %w", err)
	}
	return signupID, nil
}

func (m *Store) addCrew(ctx context.Context, tx *sqlx.Tx, s clapper.Signup) error {
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO event.crews(signup_id, position_id, locked, ordering)
		VALUES ($1, $2, $3 $4);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement to insert crew: %w", err)
	}
	for _, position := range s.Crew {
		_, err = stmt.ExecContext(ctx,
			s.SignupID, position.PositionID, position.Locked, position.Ordering)
		if err != nil {
			return fmt.Errorf("failed to insert crew for signup sheet: %w", err)
		}
	}
	return nil
}

// Update will update an existing signup sheet
func (m *Store) Update(ctx context.Context, s clapper.Signup) error {
	err := utils.Transact(m.db, func(tx *sqlx.Tx) error {
		// We'll update the metadata first
		err := m.updateMeta(ctx, tx, s)
		if err != nil {
			return err
		}
		// Then go through all the roles on the sheet
		err = m.updateCrew(ctx, tx, s)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update signup sheet: %w", err)
	}
	return nil
}

func (m *Store) updateMeta(ctx context.Context, tx *sqlx.Tx, s clapper.Signup) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE event.signups
		SET title = $1,
		SET description = $2
		SET unlock_date = $3
		SET arrival_time = $4
		SET start_time = $5
		SET end_time = $6
		WHERE signup_id = $7`, s.Title, s.Description, s.UnlockDate, s.ArrivalTime, s.StartTime, s.EndTime)
	if err != nil {
		return fmt.Errorf("failed to update signup meta: %w", err)
	}
	return nil
}

func (m *Store) updateCrew(ctx context.Context, tx *sqlx.Tx, s clapper.Signup) error {
	// Remove the previous crew
	stmt, err := tx.PrepareContext(ctx,
		`DELETE FROM events.crews WHERE crew_id = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare the update delete for roles: %w", err)
	}
	for _, pos := range s.Crew {
		_, err = stmt.ExecContext(ctx, pos.CrewID)
		if err != nil {
			return fmt.Errorf("failed to exec deleting crew from signup sheet: %w", err)
		}
	}
	// Add the updated crew
	err = m.addCrew(ctx, tx, s)
	if err != nil {
		return err
	}
	return nil
}
