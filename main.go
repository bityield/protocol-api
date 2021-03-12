package main

import (
	"net/http"

	"github.com/bityield/bityield-api/backend"
	"github.com/bityield/bityield-api/controllers"
	v1 "github.com/bityield/bityield-api/controllers/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	ginlogrus "github.com/toorop/gin-logrus"
)

// DatabaseMiddleware for gin to pass DB context around
func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("database", db)
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

func main() {
	// Initalize a new client, the base entrpy point to the application code
	b, e := backend.NewBackend()
	if e != nil {
		panic(e)
	}

	// Database connect, defer close
	defer b.R.D.Close()

	// Set the initial API instance
	r := gin.Default()

	// Use Middleware to pass around the db connection
	r.Use(ginlogrus.Logger(b.L), gin.Recovery())
	r.Use(DatabaseMiddleware(b.R.D))
	r.Use(RedisMiddleware(b.R.R))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"bityield": "Welcome to the Bityield API. Visit https://bityield.finance/developers/api for details about this API.",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("OK"))
	})

	r.StaticFile("/v1/indexes/kovan", "./assets/indexes/kovan/index.json")
	r.StaticFile("/v1/indexes/ropsten", "./assets/indexes/ropsten/index.json")

	r.GET("/ping", controllers.Ping)

	// API Methods and endpoints
	r.GET("/v1/historicals/:symbol", v1.GetHistoricals)

	// Funds endpoints
	r.GET("/funds", controllers.FindFunds)
	r.GET("/funds/:id", controllers.FindFund)

	r.Run((":" + b.C.GetString("port")))
}
