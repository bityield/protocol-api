package main

import (
	"bytes"
	"log"
	"net/http"
	"os"

	v1 "github.com/bityield/bityield-api/controllers/v1"
	"github.com/bityield/bityield-api/infra/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

// DatabaseMiddleware for gin to pass DB context around
func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("conn", db)
		c.Next()
	}
}

// RedisMiddleware for gin to pass DB context around
func RedisMiddleware(db *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", db)
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

	// Set the initial API instance
	r := gin.Default()

	// Use Middleware to pass around the db connection
	r.Use(gin.Logger())

	// Redis connection
	rd := database.ConnectRedis("localhost:6379")
	r.Use(RedisMiddleware(rd))

	// Database connection
	db := database.ConnectDatabase()
	defer db.Close()
	r.Use(DatabaseMiddleware(db))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"bityield": "Welcome to the Bityield API. Visit https://bityield.finance/developers/api for details about this API."})
	})

	r.StaticFile("/v1/indexes/kovan", "./assets/indexes/kovan/index.json")
	r.StaticFile("/v1/indexes/ropsten", "./assets/indexes/ropsten/index.json")

	// Heroku function
	r.GET("/repeat", repeatHandler(5))

	// API Methods and endpoints
	r.GET("/v1/historicals/:symbol", v1.GetHistoricals)

	// r.GET("/funds", controllers.FindFunds)
	// r.GET("/funds/:id", controllers.FindFund)

	// r.POST("/funds", controllers.CreateFund)
	// r.PATCH("/funds/:id", controllers.UpdateFund)

	r.Run((":" + port))
}
