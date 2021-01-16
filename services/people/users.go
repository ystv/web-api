package people

import (
	"context"
	"fmt"
)

// ListAll returns all users
// It doesn't return the full User object
// Returns user_id, avatar, nickname, first_name, last_name
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (m *Store) ListAll(ctx context.Context) (*[]User, error) {
	u := []User{}
	err := m.db.SelectContext(ctx, &u,
		`SELECT user_id, avatar, nickname, first_name, last_name
		FROM people.users;`)
	if err != nil {
		return nil, fmt.Errorf("fialed to list all users: %w", err)
	}
	return nil, nil
}

// ListRole returns all users who have a certain role
// It doesn't return the full User object
// Returns user_id, avatar, nickname, first_name, last_name
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (m *Store) ListRole(ctx context.Context, roleID int) (*[]User, error) {
	u := []User{}
	err := m.db.SelectContext(ctx, &u,
		`SELECT u.user_id, avatar, nickname, first_name, last_name
		FROM people.users u
		INNER JOIN people.role_members rm ON u.user_id = rm.user_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list role users: %w", err)
	}
	return nil, nil
}
