// Package campus, this is a collection of useful endpoints relating to what campus
// will hopefully offer: current term, are you on a campus net, are you
// on the ystv net
package campus

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	Repo interface {
		GetAcademicYear(ctx context.Context, t time.Time) (AcademicYear, error)
		GetTeachingPeriod(ctx context.Context, t time.Time) (TeachingPeriod, error)
		GetWeek(ctx context.Context, t time.Time) (Week, error)
		GetCurrentAcademicYear(ctx context.Context) (AcademicYear, error)
		GetCurrentTeachingPeriod(ctx context.Context) (TeachingPeriod, error)
		GetCurrentWeek(ctx context.Context) (Week, error)
	}

	// AcademicYear represents the academic year and the teaching cycle
	AcademicYear struct {
		Year          int              `db:"year" json:"year"`
		TeachingCycle []TeachingPeriod `json:"teachingCycle"`
	}

	// TeachingPeriod represents an academic time period either a term or semester
	TeachingPeriod struct {
		TeachingPeriodID int       `db:"teaching_period_id" json:"teachingPeriodID"`
		Year             int       `db:"year" json:"year"`
		Name             string    `db:"name" json:"name"` // autumn / spring / summer
		Start            time.Time `db:"start" json:"start"`
		Finish           time.Time `db:"finish" json:"finish"`
	}

	// Week is a normal week plus the number since
	// the start of a teaching period
	Week struct {
		TeachingPeriod TeachingPeriod `json:"teachingPeriod"`
		WeekNo         int            `json:"weekNo"`
	}

	Campuser struct {
		db *sqlx.DB
	}
)

var (
	ErrNoAcademicYearFound   = errors.New("no academic year found")
	ErrNoTeachingPeriodFound = errors.New("no teaching period found")
	ErrNoWeekFound           = errors.New("no week found")
)

func NewCampuser(db *sqlx.DB) Repo {
	return &Campuser{db: db}
}
