package misc

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	// QuoteRepo defines all quote interactions
	QuoteRepo interface {
		ListQuotes(ctx context.Context, amount, page int) (QuotePage, error)
		NewQuote(ctx context.Context, q Quote) error
		UpdateQuote(ctx context.Context, q Quote) error
		DeleteQuote(ctx context.Context, quoteID int) error
	}

	// WebcamRepo represents all webcam interactions
	WebcamRepo interface {
		ListWebcams(ctx context.Context, permissionIDs []int) ([]Webcam, error)
		GetWebcam(ctx context.Context, cameraID int, permissionIDs []int) (Webcam, error)
	}
)

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}
