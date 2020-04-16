package main

//go:generate ./sqlboiler --wipe psql

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/ystv/web-api/controllers"
	"github.com/ystv/web-api/utils"
)

// Get root
func rootPage(c *gin.Context) {
	c.String(200, "ystv-api - Speed and power")
}

func main() {

	r := gin.Default()

	db := utils.InitDB()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Initilise routes
	// Index
	r.GET("/", rootPage)
	// Tables
	tables := r.Group("/tables")
	{
		// Videos
		tables.GET("/videos", controllers.FindVideos)
		tables.GET("/videos/:ID", controllers.FindVideo)
		tables.POST("/videos", controllers.CreateVideo)
		// Quotes
		tables.GET("/quotes", controllers.FindQuotes)
		tables.GET("/quotes/:ID", controllers.FindQuote)
	}

	r.Run()
}
