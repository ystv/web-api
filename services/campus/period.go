package campus

import (
	"context"
	"database/sql"
	"errors"

	"time"
)

type (
	// TeachingPeriod represents an academic time period either a term or semester
	TeachingPeriod struct {
		TeachingPeriodID int       `db:"teaching_period_id" json:"teachingPeriodID"`
		Year             int       `db:"year" json:"year"`
		Name             string    `db:"name" json:"name"` // autumn / spring / summer
		Start            time.Time `db:"start" json:"start"`
		Finish           time.Time `db:"finish" json:"finish"`
	}

	// Week is a normal week plus the number since
	// the start of teaching period
	Week struct {
		TeachingPeriod TeachingPeriod `json:"teachingPeriod"`
		WeekNo         int            `json:"weekNo"`
	}
)

var (
	ErrNoTeachingPeriodFound = errors.New("no teaching period found")
	ErrNoWeekFound           = errors.New("no week found")
)

// GetTeachingPeriod retrives an academic term for a given time
func (c *Campuser) GetTeachingPeriod(ctx context.Context, t time.Time) (TeachingPeriod, error) {
	tp := TeachingPeriod{}
	err := c.db.GetContext(ctx, &tp, `
		  SELECT teaching_period_id, year, name, start, finish
		  FROM misc.teaching_periods
		  WHERE $1 BETWEEN start AND finish;`, t)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TeachingPeriod{}, ErrNoTeachingPeriodFound
		}
		return TeachingPeriod{}, err
	}
	return tp, nil
}

// GetWeek retrives the week and it's term
func (c *Campuser) GetWeek(ctx context.Context, t time.Time) (Week, error) {
	w := Week{}
	var err error
	w.TeachingPeriod, err = c.GetTeachingPeriod(ctx, t)
	if err != nil {
		if errors.Is(err, ErrNoTeachingPeriodFound) {
			return Week{}, ErrNoWeekFound
		}
		return Week{}, err
	}

	// TODO: Need to convert time from what's given to
	// that times Monday and return a better week no.
	w.WeekNo = (int(t.Sub(w.TeachingPeriod.Start).Hours()) / 24 / 7) + 1

	return w, nil
}

// GetCurrentTeachingPeriod returns the teaching period as of the current time
func (c *Campuser) GetCurrentTeachingPeriod(ctx context.Context) (TeachingPeriod, error) {
	return c.GetTeachingPeriod(ctx, time.Now())
}

// GetCurrentWeek returns the week of the current time
func (c *Campuser) GetCurrentWeek(ctx context.Context) (Week, error) {
	return c.GetWeek(ctx, time.Now())
}
