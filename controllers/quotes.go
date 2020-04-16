package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ystv/web-api/models"
)

// FindQuotes Returns all quotes
func FindQuotes(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	ctx := context.Background()
	quotes, _ := models.Quotes().All(ctx, db)
	c.JSON(http.StatusOK, quotes)
}

// FindQuote checks videos table by ID
func FindQuote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Number": "pls"})
	}

	db := c.MustGet("db").(*sql.DB)
	v := &models.Quote{ID: id}

	b, err := models.FindQuote(context.Background(), db, v.ID)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, b)
}
