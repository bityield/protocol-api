package main

import (
	"net/http"
	"time"

	"github.com/bityield/protocol-api/backend"
	v1 "github.com/bityield/protocol-api/controllers/v1"
	"github.com/bityield/protocol-api/interfaces/scrapers/coinmarketcap"
	"github.com/gin-contrib/cors"
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

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// Initalize a new client, the base entrpy point to the application code
	b, e := backend.NewBackend(true, true)
	if e != nil {
		panic(e)
	}

	// Worker jobs
	go func() {
		for now := range time.Tick(time.Second * 43200) {
			b.L.Infoln("Executing scrape job at:", now)
			coinmarketcap.Execute(b)
		}
	}()

	// Database connect, defer close
	defer b.R.D.Close()

	// Set the initial API instance
	r := gin.Default()

	// Enable and/or set cors
	cf := cors.DefaultConfig()
	cf.AllowAllOrigins = true
	cf.AllowCredentials = true
	cf.AddAllowHeaders("authorization")
	r.Use(cors.New(cf))

	// r.Use(cors.Default())
	r.Use(CORSMiddleware())

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

	// API Methods and endpoints
	r.GET("/v1/historicals/:symbol", v1.GetHistoricals)

	// Funds endpoints
	// r.GET("/funds", controllers.FindFunds)
	// r.GET("/funds/:id", controllers.FindFund)

	coinmarketcap.Execute(b)

	r.Run((":" + b.C.GetString("port")))
}
