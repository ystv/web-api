package position

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/clapper"
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
var _ clapper.PositionRepo = &Store{}

// List returns all positions
func (m *Store) List(ctx context.Context) (*[]clapper.Position, error) {
	var p []clapper.Position
	err := m.db.SelectContext(ctx, &p,
		`SELECT position_id, name, description, admin, permission_id
		FROM event.positions;`)
	if err != nil {
		err = fmt.Errorf("failed to list positions: %w", err)
		return nil, err
	}
	return &p, nil
}

// New creates a position
func (m *Store) New(ctx context.Context, p *clapper.Position) (int, error) {
	positionID := 0
	err := m.db.QueryRowContext(ctx,
		`INSERT INTO event.positions (name, description, admin, permission_id)
		VALUES ($1, $2, $3, $4) RETURNING position_id;`,
		&p.Name, &p.Description, &p.Admin, &p.PermissionID).Scan(&positionID)
	if err != nil {
		err = fmt.Errorf("failed to insert new position: %w", err)
	}
	return positionID, err
}

// Update a position, uses the ID from the token
func (m *Store) Update(ctx context.Context, p *clapper.Position) error {
	_, err := m.db.ExecContext(ctx,
		`UPDATE event.positions
		SET name=$1, description=$2, admin=$3, permission_id=$4
		WHERE position_id = $5`,
		&p.Name, &p.Description, &p.Admin, &p.PermissionID, &p.PositionID)
	if err != nil {
		err = fmt.Errorf("failed to update position: %w", err)
	}
	return err
}

// Delete a position by position ID
func (m *Store) Delete(ctx context.Context, positionID int) error {
	_, err := m.db.Exec(`DELETE FROM event.positions
						WHERE position_id = $1`, positionID)
	if err != nil {
		err = fmt.Errorf("failed to delete position: %w", err)
	}
	return err
}
