package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// GetConn ...
func GetConn(c *gin.Context) (*gorm.DB, error) {
	db, err := c.Keys["conn"].(*gorm.DB)
	if !err {
		return nil, errors.New("could not get 'conn' context connection from gin.Context")
	}

	return db, nil
}

// GetRedis ...
func GetRedis(c *gin.Context) (*redis.Client, error) {
	db, err := c.Keys["redis"].(*redis.Client)
	if !err {
		return nil, errors.New("could not get 'redis' context connection from gin.Context")
	}

	return db, nil
}
