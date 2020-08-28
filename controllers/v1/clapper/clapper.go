package clapper

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/clapper"
	"github.com/ystv/web-api/services/clapper/event"
	"github.com/ystv/web-api/services/clapper/position"
	"github.com/ystv/web-api/services/clapper/signup"
)

// Repos encapsulates the dependency
type Repos struct {
	event    clapper.EventRepo
	signup   clapper.SignupRepo
	position clapper.PositionRepo
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{event.NewStore(db), signup.NewStore(db), position.NewStore(db)}
}
