package public

import (
	"github.com/jackc/pgx"
	"github.com/ystv/web-api/services/public"
)

// Repos encapsulates the dependency
type Repos struct {
	public *public.Store
}

// NewRepos creates our data store
func NewRepos(db *pgx.Conn) *Repos {
	return &Repos{public.NewStore(db)}
}
