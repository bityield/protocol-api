package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"

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

// repeatHandler for Heroku
func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	repeat, err := strconv.Atoi(os.Getenv("REPEAT"))
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	// Set the initial API instance
	r := gin.Default()

	// Database connection
	db := database.ConnectDatabase()
	defer db.Close()

	// Use Middleware to pass around the db connection
	r.Use(gin.Logger())
	r.Use(DatabaseMiddleware(db))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"bityield": "Welcome to the Bityield API. Visit https://bityield.finance/developers/api for details about this API."})
	})

	r.StaticFile("/v1/indexes/kovan", "./assets/indexes/kovan/index.json")

	// Heroku function
	r.GET("/repeat", repeatHandler(repeat))

	// r.GET("/funds", controllers.FindFunds)
	// r.GET("/funds/:id", controllers.FindFund)

	// r.POST("/funds", controllers.CreateFund)
	// r.PATCH("/funds/:id", controllers.UpdateFund)

	r.Run((":" + port))
}
