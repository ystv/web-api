package clapper

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/clapper"
	"github.com/ystv/web-api/services/clapper/event"
	"github.com/ystv/web-api/services/clapper/position"
)

// Repos encapsulates the dependency
type Repos struct {
	event    clapper.EventRepo
	position clapper.PositionRepo
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{event.NewStore(db), position.NewStore(db)}
}
