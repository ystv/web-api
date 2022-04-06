package clapper

import (
	"github.com/jackc/pgx/v4"

	"github.com/ystv/web-api/services/clapper"
	"github.com/ystv/web-api/services/clapper/crew"
	"github.com/ystv/web-api/services/clapper/event"
	"github.com/ystv/web-api/services/clapper/position"
	"github.com/ystv/web-api/services/clapper/signup"
	"github.com/ystv/web-api/utils"
)

// Repos encapsulates the dependency
type Repos struct {
	access   *utils.Accesser
	crew     clapper.CrewRepo
	event    clapper.EventRepo
	signup   clapper.SignupRepo
	position clapper.PositionRepo
}

// NewRepos creates our data store
func NewRepos(db *pgx.Conn, access *utils.Accesser) *Repos {
	return &Repos{
		access,
		crew.NewStore(db),
		event.NewStore(db),
		signup.NewStore(db),
		position.NewStore(db),
	}
}
