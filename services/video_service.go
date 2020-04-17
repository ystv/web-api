package services

import (
	"context"

	"github.com/ystv/web-api/models"
	"github.com/ystv/web-api/utils"
)

// VideoCreate new video create
// func VideoCreate(displayname string) (uint64, error) {
// 	v := &models.Video{DisplayName: displayname}
// }

// VideoList list video
func VideoList() (models.VideoSlice, error) {
	ctx := context.Background()
	v, err := models.Videos().All(ctx, utils.DB)
	return v, err
}
