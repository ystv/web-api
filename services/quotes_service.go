package services

import (
	"context"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/ystv/web-api/models"
	"github.com/ystv/web-api/utils"
)

// QuoteCreate new quote create
func QuoteCreate(q *models.Quote) (int, error) {
	ctx := context.Background()
	err := q.Insert(ctx, utils.DB, boil.Infer())
	return q.ID, err
}

// QuoteList list quote
func QuoteList() (models.QuoteSlice, error) {
	ctx := context.Background()
	q, err := models.Quotes().All(ctx, utils.DB)
	return q, err
}

// QuoteFind find quote
func QuoteFind(id int) (*models.Quote, error) {
	ctx := context.Background()
	q, err := models.FindQuote(ctx, utils.DB, id)
	return q, err
}

// QuoteUpdate update quote
func QuoteUpdate(oldQuote *models.Quote, newQuote *models.Quote) error {
	oldQuote.Quote = newQuote.Quote
	oldQuote.Description = newQuote.Description
	oldQuote.PostedBy = newQuote.PostedBy
	ctx := context.Background()
	_, err := oldQuote.Update(ctx, utils.DB, boil.Infer())
	return err
}

// QuoteDelete
func QuoteDelete(id int) (*models.Quote, error) {
	q, err := QuoteFind(id)
	ctx := context.Background()
	_, err = q.Delete(ctx, utils.DB)
	return q, err
}
