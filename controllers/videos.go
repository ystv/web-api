package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ystv/web-api/models"
)

// FindVideos Returns all videos
func FindVideos(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	ctx := context.Background()
	videos, _ := models.Videos().All(ctx, db)
	c.JSON(http.StatusOK, gin.H{"data": videos})
}

// FindVideo checks videos table by ID
func FindVideo(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Number": "pls"})
	}

	db := c.MustGet("db").(*sql.DB)
	v := &models.Video{ID: id}

	found, err := models.FindVideo(context.Background(), db, v.ID)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": found})
}
