package misc

import (
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/misc"
	"github.com/ystv/web-api/utils"
)

// Repos stores our dependencies
type Repos struct {
	misc   misc.Repos
	access utils.Repo
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB, access utils.Repo) *Repos {
	return &Repos{
		misc:   misc.NewStore(db),
		access: access,
	}
}
