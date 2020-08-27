package public

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/public"
)

// Repos encapsulates the dependency
type Repos struct {
	public *public.Store
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{public.NewStore(db)}
}
