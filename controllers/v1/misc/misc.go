package misc

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/misc"
	"github.com/ystv/web-api/utils"
)

// Repos stores our dependencies
type Repos struct {
	misc   *misc.Store
	access *utils.Accesser
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB, access *utils.Accesser) *Repos {
	return &Repos{
		misc:   misc.NewStore(db),
		access: access,
	}
}
