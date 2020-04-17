package services

import (
	"context"

	"github.com/ystv/web-api/models"
	"github.com/ystv/web-api/utils"
)

// VideoBoxCreate new quote create
// func quoteCreate(displayname string) (uint64, error) {
// 	v := &models.Quote{DisplayName: displayname}
// }

// VideoBoxList list quote
func VideoBoxList() (models.VideoBoxSlice, error) {
	ctx := context.Background()
	b, err := models.VideoBoxes().All(ctx, utils.DB)
	return b, err
}

// VideoBoxFind find quote
func VideoBoxFind(id int) (*models.VideoBox, error) {
	ctx := context.Background()
	v, err := models.FindVideoBox(ctx, utils.DB, id)
	return v, err
}
