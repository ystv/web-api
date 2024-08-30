package people

import (
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

type (
	// Repo stores our dependencies
	Repo struct {
		people *people.Store
		access *utils.Accesser
	}
)

// NewRepo creates our data store
func NewRepo(db *sqlx.DB, access *utils.Accesser) *Repo {
	return &Repo{
		people: people.NewStore(db),
		access: access,
	}
}
