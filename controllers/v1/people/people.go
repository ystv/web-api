package people

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/people"
)

// Repo stores our dependencies
type Repo struct {
	people *people.Store
}

// NewRepo creates our data store
func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{people.NewStore(db)}
}
