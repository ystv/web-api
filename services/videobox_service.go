package services

import (
	"context"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/ystv/web-api/models"
	"github.com/ystv/web-api/utils"
)

// VideoBoxCreate new videobox create
func VideoBoxCreate(vb *models.VideoBox) (int, error) {
	ctx := context.Background()
	err := vb.Insert(ctx, utils.DB, boil.Infer())
	return vb.ID, err
}

// VideoBoxList list videoboxes
func VideoBoxList() (models.VideoBoxSlice, error) {
	ctx := context.Background()
	vb, err := models.VideoBoxes().All(ctx, utils.DB)
	return vb, err
}

// VideoBoxFind find videobox
func VideoBoxFind(id int) (*models.VideoBox, error) {
	ctx := context.Background()
	vb, err := models.FindVideoBox(ctx, utils.DB, id)
	return vb, err
}

// VideoBoxUpdate update quote
func VideoBoxUpdate(oldVideoBox *models.VideoBox, newVideoBox *models.VideoBox) error {
	// TODO Rewrite updating
	oldVideoBox.DisplayName = newVideoBox.DisplayName
	oldVideoBox.Description = newVideoBox.Description
	oldVideoBox.Image = newVideoBox.Image
	oldVideoBox.IsEnabled = newVideoBox.IsEnabled
	oldVideoBox.IsProduction = newVideoBox.IsProduction
	oldVideoBox.IsPublic = newVideoBox.IsPublic
	oldVideoBox.IsVisibleInLatestVideos = newVideoBox.IsVisibleInLatestVideos
	ctx := context.Background()
	_, err := oldVideoBox.Update(ctx, utils.DB, boil.Infer())
	return err
}

// VideoBoxDelete
func VideoBoxDelete(id int) (*models.VideoBox, error) {
	vb, err := VideoBoxFind(id)
	ctx := context.Background()
	_, err = vb.Delete(ctx, utils.DB)
	return vb, err
}
