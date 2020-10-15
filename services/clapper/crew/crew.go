package crew

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
var _ clapper.CrewRepo = &Store{}

// Get returns a crew position object
func (m *Store) Get(ctx context.Context, crewID int) (*clapper.CrewPosition, error) {
	cp := clapper.CrewPosition{}
	err := m.db.GetContext(ctx, &cp,
		`SELECT crew_id, signup_id, user_id, locked, admin, permission_id
		FROM event.crews
		WHERE crew_id = $1;`, crewID)
	if err != nil {
		err = fmt.Errorf("failed to get crew from crewID: %w", err)
		return nil, err
	}
	return &cp, nil
}

// UpdateUser Updates the user field for the specified crew ID to the specified user ID
func (m *Store) UpdateUser(ctx context.Context, crewID, userID int) error {
	return utils.Transact(m.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareContext(ctx,
			`UPDATE event.crews
			SET user_id = $1
			WHERE crew_id = $2;`)
		if err != nil {
			err = fmt.Errorf("failed to prepare statement to update crew: %w", err)
			return err
		}
		_, err = stmt.ExecContext(ctx, userID, crewID)
		if err != nil {
			err = fmt.Errorf("failed to execute statement on crew update: %w", err)
			return err
		}
		return nil
	})
}

// DeleteUser clears the user ID from the crew ID object
func (m *Store) DeleteUser(ctx context.Context, crewID int) error {
	return utils.Transact(m.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareContext(ctx,
			`UPDATE event.crews
			SET user_id = NULL
			WHERE crew_id = $2;`)
		if err != nil {
			err = fmt.Errorf("failed to prepare statement to delete crew user: %w", err)
			return err
		}
		_, err = stmt.ExecContext(ctx, crewID)
		if err != nil {
			err = fmt.Errorf("failed to execute statement on crew delete user: %w", err)
			return err
		}
		return nil
	})
}
