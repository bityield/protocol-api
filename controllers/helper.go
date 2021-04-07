package controllers

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

var (
	CTX = context.Background()

	HOR_INTERVAL = 3600
	DAY_INTERVAL = 86400
	WEK_INTERVAL = 604800
	MNT_INTERVAL = 2592000
	YER_INTERVAL = 31536000

	INT_HOR = "h"
	INT_DAY = "d"
	INT_WEK = "w"
	INT_MNT = "m"
	INT_YER = "y"

	Intervals = map[string]int{
		INT_HOR: HOR_INTERVAL,
		INT_DAY: DAY_INTERVAL,
		INT_WEK: WEK_INTERVAL,
		INT_MNT: MNT_INTERVAL,
		INT_YER: YER_INTERVAL,
	}
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
