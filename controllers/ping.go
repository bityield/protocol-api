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

	db.Set(ctx, "ping", "pong", 0)

	c.JSON(http.StatusOK, gin.H{
		"ping": db.Get(ctx, "ping").Val(),
	})
}
