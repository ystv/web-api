package misc

import (
	"github.com/jackc/pgx"
	"github.com/ystv/web-api/services/misc"
	"github.com/ystv/web-api/utils"
)

// Repos stores our dependencies
type Repos struct {
	misc   *misc.Store
	access *utils.Accesser
}

// NewRepos creates our data store
func NewRepos(db *pgx.Conn, access *utils.Accesser) *Repos {
	return &Repos{
		misc:   misc.NewStore(db),
		access: access,
	}
}
