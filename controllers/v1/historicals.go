package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bityield/protocol-api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// GetHistoricals ...
func GetHistoricals(c *gin.Context) {
	db, err := controllers.GetRedis(c)
	if err != nil {
		panic(err)
	}

	var interval, min, max string

	interval = c.Query("interval")
	if interval == "" {
		interval = "day"
	}

	min = c.Query("min")
	if min == "" {
		min = "-inf"
	}

	max = c.Query("max")
	if max == "" {
		max = "+inf"
	}

	symbol := strings.ToLower(c.Param("symbol"))

	vals, err := db.ZRangeByScoreWithScores(ctx, symbol, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
	}).Result()

	if err != nil {
		panic(err)
	}

	ohlcv := []map[string]interface{}{}
	dates := []int{}
	prices := []string{}

	for _, v := range vals {
		var member map[string]interface{}
		json.Unmarshal([]byte(v.Member.(string)), &member)

		ohlcv = append(ohlcv, member)
		dates = append(dates, int(v.Score))
		prices = append(prices, member["open"].(string))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ohlcv,
	})
}
