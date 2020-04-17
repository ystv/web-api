package services

import (
	"github.com/ystv/web-api/models"
)

// VideoCreate new video create
func VideoCreate(displayname string) (uint64, error) {
	v := &models.Video{DisplayName: displayname}
}

// VideoList list video
func VideoList() ([]*models.VideoList, error) {
	v := &models.Video{}
	return v.List()
}
