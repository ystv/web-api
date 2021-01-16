package people

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// UserRepo defines all user interactions
type UserRepo interface {
	Get(ctx context.Context, userID int) (*User, error)
	GetFull(ctx context.Context, userID int) (*UserFull, error)
	ListAll(ctx context.Context) (*[]User, error)
	ListRole(ctx context.Context, roleID int) (*[]User, error)
}

// PermissionRepo defines all permission interactions
type PermissionRepo interface {
}

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// Here for validation to ensure we are meeting the interface
var _ UserRepo = &Store{}
