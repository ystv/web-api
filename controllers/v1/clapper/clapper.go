package clapper

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/clapper"
	"github.com/ystv/web-api/services/clapper/crew"
	"github.com/ystv/web-api/services/clapper/event"
	"github.com/ystv/web-api/services/clapper/position"
	"github.com/ystv/web-api/services/clapper/signup"
	"github.com/ystv/web-api/utils"
)

// Repos encapsulates the dependency
type (
	Repos interface {
		CrewRepo
		EventRepo
		PositionRepo
		SignupRepo
	}

	CrewRepo interface {
		SetCrew(c echo.Context) error
		ResetCrew(c echo.Context) error
		NewCrew(c echo.Context) error
		DeleteCrew(c echo.Context) error
	}

	EventRepo interface {
		ListMonth(c echo.Context) error
		GetEvent(c echo.Context) error
		NewEvent(c echo.Context) error
		UpdateEvent(c echo.Context) error
		DeleteEvent(c echo.Context) error
	}

	PositionRepo interface {
		ListPositions(c echo.Context) error
		NewPosition(c echo.Context) error
		UpdatePosition(c echo.Context) error
		DeletePosition(c echo.Context) error
	}

	SignupRepo interface {
		NewSignup(c echo.Context) error
		UpdateSignup(c echo.Context) error
		DeleteSignup(c echo.Context) error
	}

	Store struct {
		access   utils.Repo
		crew     clapper.CrewRepo
		event    clapper.EventRepo
		signup   clapper.SignupRepo
		position clapper.PositionRepo
	}
)

// NewRepos creates our data store
func NewRepos(db *sqlx.DB, access utils.Repo) Repos {
	return &Store{
		access,
		crew.NewStore(db),
		event.NewStore(db),
		signup.NewStore(db),
		position.NewStore(db),
	}
}
