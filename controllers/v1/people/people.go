package people

import (
	"github.com/jackc/pgx"
	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

// Repo stores our dependencies
type Repo struct {
	people *people.Store
	access *utils.Accesser
}

// NewRepo creates our data store
func NewRepo(db *pgx.Conn, access *utils.Accesser) *Repo {
	return &Repo{
		people: people.NewStore(db),
		access: access,
	}
}
