package controllers

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/ystv/web-api/models"
)

// FindVideoFiles Returns all video files
func FindVideoFiles(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	ctx := context.Background()
	vf, _ := models.VideoFiles().All(ctx, db)
	c.JSON(200, vf)
}

// FindVideoFile checks video_files table by ID
func FindVideoFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.JSON(400, gin.H{"Number": "pls"})
	}

	db := c.MustGet("db").(*sql.DB)
	vf := &models.VideoFile{ID: id}

	b, err := models.FindVideoFile(context.Background(), db, vf.ID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, b)
}

// CreateVideoFile creates a video file
func CreateVideoFile(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Validate input
	v := &models.VideoFile{}
	if err := c.ShouldBindJSON(v); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	err := v.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		c.String(400, "Invalid request")
	}
}
