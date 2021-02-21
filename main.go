package main

import (
	"net/http"

	"github.com/bityield/bityield-api/controllers"
	"github.com/bityield/bityield-api/infra/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DatabaseMiddleware for gin to pass DB context around
func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("conn", db)
		c.Next()
	}
}

func main() {
	// Set the initial API instance
	r := gin.Default()

	// Database connection
	db := database.ConnectDatabase()
	defer db.Close()

	// Use Middleware to pass around the db connection
	r.Use(DatabaseMiddleware(db))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"bityield": "Welcome to the Bityield API. Visit https://bityield.finance/developers/api for details about this API."})
	})

	r.GET("/funds", controllers.FindFunds)
	r.GET("/funds/:id", controllers.FindFund)

	r.POST("/funds", controllers.CreateFund)
	r.PATCH("/funds/:id", controllers.UpdateFund)

	r.Run()
}
