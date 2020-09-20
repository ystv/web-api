package misc

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/misc"
)

// Repos stores our dependencies
type Repos struct {
	misc *misc.Store
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{misc.NewStore(db)}
}
