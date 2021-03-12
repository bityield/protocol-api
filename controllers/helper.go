package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

// GetDatabase ...
func GetDatabase(c *gin.Context) (*gorm.DB, error) {
	db, err := c.Keys["database"].(*gorm.DB)
	if !err {
		return nil, errors.New("could not get 'database' context connection from gin.Context")
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
