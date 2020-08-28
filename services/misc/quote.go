package misc

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// Quote is an individual quote
	Quote struct {
		QuoteID     int    `db:"quote_id" json:"id"`
		Quote       string `db:"quote" json:"quote"`
		Description string `db:"description" json:"description"`
		CreatedBy   int    `db:"created_by" json:"createdBy"`
	}
	// QuotePage is a group of quotes including the last page idnex
	QuotePage struct {
		Quotes        []Quote
		LastPageIndex int
	}
)

// Repo defines all misc interactions
type Repo interface {
	ListQuotes(ctx context.Context, amount, page int) (QuotePage, error)
	NewQuote(ctx context.Context, q Quote) error
	UpdateQuote(ctx context.Context, q Quote) error
	DeleteQuote(ctx context.Context, quoteID int) error
}

// Here for validation to ensure we are meeting the interface
var _ Repo = &Store{}

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// ListQuotes returns a section of quotes
func (s *Store) ListQuotes(ctx context.Context, amount, page int) (QuotePage, error) {
	q := QuotePage{LastPageIndex: 20}
	err := s.db.SelectContext(ctx, &q.Quotes,
		`SELECT quote_id, quote, description, created_by
		FROM misc.quotes
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2;`, amount, page)
	return q, err
}

// NewQuote creates a new quote
func (s *Store) NewQuote(ctx context.Context, q Quote) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO misc.quotes(quote, description, created_at, created_by)
		VALUES ($1, $2, $3, $4);`, q.Quote, q.Description, q.CreatedBy, time.Now())
	return err
}

// UpdateQuote updates a quote
func (s *Store) UpdateQuote(ctx context.Context, q Quote) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE misc.quotes SET quote = $1, description = $2
	WHERE quote_id = $3;`, q.Quote, q.Description, q.QuoteID)
	return err
}

// DeleteQuote deletes a quote
func (s *Store) DeleteQuote(ctx context.Context, quoteID int) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM misc.quotes WHERE quote_id = $1;`, quoteID)
	return err
}
