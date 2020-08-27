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
		(event_id, title, description, unlock_date, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING signup_id;`,
			eventID, s.Title, s.Description, s.UnlockDate, s.StartTime, s.EndTime).Scan(&signupID)
		if err != nil {
			return fmt.Errorf("failed to insert new signup shee metat: %w", err)
		}
		// Check if positions have been added
		// TODO I'm not too sure on using the signup struct for this,
		// maybe another input variable instead?
		if len(s.Crew) == 0 {
			return nil
		}
		stmt, err := tx.PrepareContext(ctx,
			`INSERT INTO event.crews(signup_id, position_id, locked, ordering)
			VALUES ($1, $2, $3 $4);`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement to insert crew: %w", err)
		}
		for _, position := range s.Crew {
			_, err = stmt.ExecContext(ctx,
				signupID, position.PositionID, position.Locked, position.Ordering)
			if err != nil {
				return fmt.Errorf("failed to insert crew for signup sheet: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert new signup sheet: %w", err)
	}
	return signupID, nil
}
