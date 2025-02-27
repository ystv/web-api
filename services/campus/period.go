package campus

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// GetAcademicYear retrieves an academic year for a given time
func (c *Campuser) GetAcademicYear(ctx context.Context, t time.Time) (AcademicYear, error) {
	ay := AcademicYear{}
	ay.Year = t.Year()

	err := c.db.SelectContext(ctx, &ay.TeachingCycle, `
	SELECT period_id, period.year, name, start, finish
	FROM misc.teaching_periods period
	INNER JOIN (
		SELECT year
		FROM misc.teaching_periods
		GROUP BY year
		HAVING $1 BETWEEN min(start) AND max(finish)
	) selected_year ON selected_year.year = period.year
	ORDER BY start;`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AcademicYear{}, ErrNoAcademicYearFound
		}
		return AcademicYear{}, err
	}
	return ay, nil
}

// GetTeachingPeriod retrieves an academic term for a given time
func (c *Campuser) GetTeachingPeriod(ctx context.Context, t time.Time) (TeachingPeriod, error) {
	tp := TeachingPeriod{}
	err := c.db.GetContext(ctx, &tp, `
		  SELECT period_id, year, name, start, finish
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

// GetWeek retrieves the week and it's term
func (c *Campuser) GetWeek(ctx context.Context, t time.Time) (Week, error) {
	w := Week{}
	var err error
	w.TeachingPeriod, err = c.GetTeachingPeriod(ctx, t)
	if err != nil {
		if errors.Is(err, ErrNoTeachingPeriodFound) {
			return Week{}, ErrNoWeekFound
		}
		return Week{}, fmt.Errorf("failed to get teaching period: %w", err)
	}

	// TODO: Need to convert time from what's given to that times Monday and return a better week no.
	w.WeekNo = (int(t.Sub(w.TeachingPeriod.Start).Hours()) / 24 / 7) + 1

	return w, nil
}

// GetCurrentAcademicYear returns the academic year as of the current time
func (c *Campuser) GetCurrentAcademicYear(ctx context.Context) (AcademicYear, error) {
	return c.GetAcademicYear(ctx, time.Now())
}

// GetCurrentTeachingPeriod returns the teaching period as of the current time
func (c *Campuser) GetCurrentTeachingPeriod(ctx context.Context) (TeachingPeriod, error) {
	return c.GetTeachingPeriod(ctx, time.Now())
}

// GetCurrentWeek returns the week of the current time
func (c *Campuser) GetCurrentWeek(ctx context.Context) (Week, error) {
	return c.GetWeek(ctx, time.Now())
}
