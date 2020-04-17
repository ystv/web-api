package tables

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/ystv/web-api/models"
)

// FindVideoBoxes Returns all videos
func FindVideoBoxes(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	ctx := context.Background()
	videos, _ := models.VideoBoxes().All(ctx, db)
	c.JSON(200, videos)
}

// FindVideoBox checks video_boxes table by ID
func FindVideoBox(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.JSON(400, gin.H{"Number": "pls"})
	}

	db := c.MustGet("db").(*sql.DB)
	vb := &models.VideoBox{ID: id}

	b, err := models.FindVideoFile(context.Background(), db, vb.ID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, b)
}

// CreateVideoBox creates a videobox
func CreateVideoBox(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Validate input
	v := &models.VideoBox{}
	if err := c.ShouldBindJSON(v); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	err := v.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		c.String(400, "Invalid request")
	}
}
