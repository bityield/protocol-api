package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// Ping ...
func Ping(c *gin.Context) {
	db, err := GetRedis(c)
	if err != nil {
		panic(err)
	}

	val, err := db.Get(ctx, "ping").Result()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"ping": val,
	})
}
