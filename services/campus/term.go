package campus

import (
	"context"
	"time"
)

type (
	// Term represents an academic term
	Term struct {
		TermID int       `db:"term_id" json:"termID"`
		Year   string    `db:"year" json:"year"`
		Term   string    `db:"term" json:"term"` // autumn / spring / summer
		Start  time.Time `db:"start" json:"start"`
		Finish time.Time `db:"finish" json:"finish"`
	}
	// Week is a normal week plus the number since
	// the start of term
	Week struct {
		Term   Term
		WeekNo int
	}
)

// GetTerm retrives an academic term for a given time
func (c *Campuser) GetTerm(ctx context.Context, t time.Time) (Term, error) {
	term := Term{}
	c.db.GetContext(ctx, &t, `
		SELECT term_id, year, term, start, finish
		FROM misc.terms
		WHERE $1 BETWEEN start AND FINISH;`, t)
	return term, nil
}

// GetWeek retrives the week and it's term
func (c *Campuser) GetWeek(ctx context.Context, t time.Time) (Week, error) {
	// TODO: This is horrible. Need to convert time from what's give to
	// that times Monday and return a better week no.
	w := Week{}
	c.db.GetContext(ctx, &w.Term, `
		SELECT term_id, year, term, start, finish
		FROM misc.terms
		WHERE $1 BETWEEN start AND FINISH;`, t)
	w.WeekNo = (int(t.Sub(w.Term.Start).Hours()) / 24 / 7) + 1
	return w, nil
}

// GetCurrentTerm returns the term as of the current time
func (c *Campuser) GetCurrentTerm(ctx context.Context) (Term, error) {
	return c.GetTerm(ctx, time.Now())
}

// GetCurrentWeek returns the week of the current time
func (c *Campuser) GetCurrentWeek(ctx context.Context) (Week, error) {
	return c.GetWeek(ctx, time.Now())
}
