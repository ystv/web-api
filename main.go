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

	r.GET("/", rootPage)
	r.GET("/videos", controllers.FindVideos)
	r.GET("/videos/:ID", controllers.FindVideo)
	r.Run()
}
