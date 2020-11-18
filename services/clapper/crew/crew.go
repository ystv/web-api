package crew

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/guregu/null.v4"

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

type crew struct {
	userID       null.Int `db:"user_id"`
	PermissionID null.Int `db:"permission_id"`
}

// Get returns a crew position object
func (m *Store) Get(ctx context.Context, crewID int) (*clapper.CrewPosition, error) {
	cp := clapper.CrewPosition{}
	err := m.db.GetContext(ctx, &cp,
		`SELECT crew_id, user_id, locked, admin, permission_id
		FROM event.crews
		INNER JOIN event.positions ON crews.position_id = positions.position_id
		WHERE crew_id = $1;`, crewID)
	if err != nil {
		err = fmt.Errorf("failed to get crew from crewID: %w", err)
		return nil, err
	}
	return &cp, nil
}

// DeleteUser clears the user ID from the crew ID object
func (m *Store) DeleteUser(ctx context.Context, crewID int) error {
	return utils.Transact(m.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareContext(ctx,
			`UPDATE event.crews
			SET user_id = NULL
			WHERE crew_id = $1;`)
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

// updateUser updates userID in a crew object
func (m *Store) updateUser(ctx context.Context, tx *sqlx.Tx, crewID, userID int) error {
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
}

// UpdateUser Updates the user field for the specified crew ID to the specified user ID
func (m *Store) UpdateUser(ctx context.Context, crewID, userID int) error {
	return utils.Transact(m.db, func(tx *sqlx.Tx) error {
		return m.updateUser(ctx, tx, crewID, userID)
	})
}

// UpdateUserAndVerify will update a crew object to the specified user ID,
// it will also perform additional checks to ensure they have enough permission
func (m *Store) UpdateUserAndVerify(ctx context.Context, eventID, crewID, userID int) error {
	err := utils.Transact(m.db, func(tx *sqlx.Tx) error {
		// we're just checking if a user has already signed up, otherwise go for it
		crew, err := m.checkSameUser(ctx, tx, crewID, userID)
		if err != nil {
			return err
		}
		if !crew.userID.Valid {
			// no-one has signed-up, check if they have permission
			if crew.PermissionID.Valid {
				// role does require permission
				err = m.checkUserRole(ctx, tx, crewID, userID)
				if err != nil {
					return err
				}
				// they do have permission, carry on with updating
				err = m.updateUser(ctx, tx, crewID, userID)
				if err != nil {
					return err
				}
			}
		}
		// they are kicking someone off, so lets check they have consent from the government (authorization)
		// check if they are an admin of the event
		err = m.checkEventAdmin(ctx, tx, eventID, userID)
		if err != nil {
			return err
		}
		// at this point, all checks are complete and they can change
		return m.updateUser(ctx, tx, crewID, userID)
	})
	if err != nil {
		return fmt.Errorf("failed to update and verify user: %w", err)
	}
	return nil
}

func (m *Store) getEventFromCrew(ctx context.Context, tx *sqlx.Tx, crewID int) error {
	return nil
}

// Misc authorization checks

func (m *Store) checkSuperUser(ctx context.Context, tx *sqlx.Tx, userID int) error {
	stmt, err := tx.PrepareContext(ctx,
		`SELECT true
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1 AND p.permission_id = 19;`)
	if err != nil {
		return fmt.Errorf("failed to prepare super user check: %w", err)
	}
	isSuperUser := null.Bool{}
	err = stmt.QueryRowContext(ctx, userID).Scan(&isSuperUser)
	if err != nil {
		return fmt.Errorf("failed to exec super user check: %w", err)
	}
	if !isSuperUser.Valid {
		return errors.New("user doesn't have super-user permission")
	}
	return nil
}

func (m *Store) checkEventAdmin(ctx context.Context, tx *sqlx.Tx, eventID, userID int) error {
	stmt, err := tx.PrepareContext(ctx,
		`SELECT bool_or(position.admin) AS has_admin
		FROM event.crews crew
		INNER JOIN event.positions position ON crew.position_id = position.position_id
		INNER JOIN event.signups signup ON crew.signup_id = signup.signup_id
		INNER JOIN event.events event ON signup.event_id = event.event_id
		WHERE event.event_id = $1 AND crew.user_id = $2;`)
	if err != nil {
		return fmt.Errorf("failed to prepare event admin check: %w", err)
	}
	hasAdmin := false
	err = stmt.QueryRowContext(ctx, eventID, userID).Scan(&hasAdmin)
	if err != nil {
		return fmt.Errorf("failed to query event admin check: %w", err)
	}
	if !hasAdmin {
		return errors.New("user doesn't have admin permission")
	}
	return nil
}

func (m *Store) checkUserRole(ctx context.Context, tx *sqlx.Tx, crewID, userID int) error {
	stmt, err := tx.PrepareContext(ctx,
		`SELECT CASE WHEN EXISTS(
						SELECT true
						FROM event.crews crew
						INNER JOIN event.positions position ON crew.position_id = position.position_id;
						INNER JOIN people.role_permissions permission ON position.permission_id = permission.permission_id;
						INNER JOIN people.role_members member ON permission.role_id = member.role_id
						WHERE crew.crew_id = $1 AND member.user_id = $2
					)
						THEN true
						ELSE false
					END AS has_permission;`)
	if err != nil {
		return fmt.Errorf("failed to prepare check permission statement: %w", err)
	}
	hasPermission := false
	err = stmt.QueryRowContext(ctx, crewID, userID).Scan(&hasPermission)
	if err != nil {
		return fmt.Errorf("failed to check permission of user for crew: %w", err)
	}
	if !hasPermission {
		return errors.New("user doesn't have role permission")
	}
	return nil
}

func (m *Store) checkSameUser(ctx context.Context, tx *sqlx.Tx, crewID, userID int) (crew, error) {
	stmt, err := tx.PrepareContext(ctx,
		`SELECT crew.user_id, position.permission_id
	FROM event.crews crew
	INNER JOIN positions position
	ON crew.position_id = position.position_id
	WHERE crew.crew_id = $1 AND crew.user_id = $2;
	`)
	crew := crew{}
	if err != nil {
		return crew, fmt.Errorf("failed to prepare same user statement: %w", err)
	}

	err = stmt.QueryRowContext(ctx, crewID, userID).Scan(&crew)
	if err != nil {
		return crew, fmt.Errorf("failed to query crew user and perm: %w", err)
	}
	return crew, nil
}
