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
func (m *Store) New(ctx context.Context, eventID int, s clapper.NewSignup) (int, error) {
	signupID := 0
	err := utils.Transact(m.db, func(tx *sqlx.Tx) error {
		err := tx.QueryRowContext(ctx, `INSERT INTO event.signups
		(event_id, title, description, unlock_date, arrival_time, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING signup_id;`,
			eventID, s.Title, s.Description, s.UnlockDate, s.ArrivalTime, s.StartTime, s.EndTime).Scan(&signupID)
		if err != nil {
			return fmt.Errorf("failed to insert new signup sheet: %w", err)
		}
		// Check if positions have been added
		// TODO I'm not too sure on using the signup struct for this,
		// maybe another input variable instead?
		if len(s.Crew) == 0 {
			return nil
		}
		err = m.addCrew(ctx, tx, signupID, s.Crew)
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

func (m *Store) addCrew(ctx context.Context, tx *sqlx.Tx, signupID int, crew []clapper.NewCrew) error {
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO event.crews(signup_id, position_id, locked, credited, ordering)
		VALUES ($1, $2, $3 $4);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement to insert crew: %w", err)
	}
	for _, position := range crew {
		// we might want to handle the signup ID if it does exist, just as an if a statement does
		_, err = stmt.ExecContext(ctx,
			signupID, position.PositionID, position.Locked, position.Credited, position.Ordering)
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
		description = $2,
		unlock_date = $3,
		arrival_time = $4,
		start_time = $5,
		end_time = $6
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
	err = m.addCrew(ctx, tx, s.SignupID, m.crewPositionToNewCrew(s.Crew))
	if err != nil {
		return err
	}
	return nil
}

func (m *Store) crewPositionToNewCrew(crewPosition []clapper.CrewPosition) []clapper.NewCrew {
	newCrew := make([]clapper.NewCrew, 0)
	for _, crew := range crewPosition {
		newCrew = append(newCrew, clapper.NewCrew{
			PositionID: crew.PositionID,
			Locked:     crew.Locked,
			Credited:   crew.Credited,
		})
	}
	return newCrew
}

// Delete will remove the signup sheet, and it's children crew
// (The database should cascade to delete children)
func (m *Store) Delete(ctx context.Context, signupID int) error {
	err := utils.Transact(m.db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx,
			`DELETE FROM events.signups WHERE signup_id = $1`, signupID)
		if err != nil {
			return fmt.Errorf("failed to exec deleting signup sheet: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete signup sheet: %w", err)
	}
	return nil
}
