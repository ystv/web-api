package misc

import (
	"context"

	"github.com/jackc/pgx"
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
		ListWebcams(ctx context.Context, permissions []string) ([]Webcam, error)
		GetWebcam(ctx context.Context, cameraID int, permissions []string) (Webcam, error)
	}
	// ListRepo represents all mailing list interactions
	//
	// Send emails to roles and other groups
	ListRepo interface {
		GetLists(ctx context.Context) ([]List, error)
		GetListsByUserID(ctx context.Context, userID int) ([]List, error)
		GetList(ctx context.Context, listID int) (List, error)
		GetSubscribers(ctx context.Context, listID int) ([]Subscriber, error)
		Subscribe(ctx context.Context, userID, listID int) error
		UnsubscribeByID(ctx context.Context, userID, listID int) error
		UnsubscribeByUUID(ctx context.Context, uuid string) error
	}
)

// Store contains our dependency
type Store struct {
	db *pgx.Conn
}

// NewStore creates a new store
func NewStore(db *pgx.Conn) *Store {
	return &Store{db: db}
}
