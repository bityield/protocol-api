package coins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bityield/protocol-api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// GetCoinHistoricals ...
func GetCoinHistoricals(c *gin.Context) {
	db, err := controllers.GetRedis(c)
	if err != nil {
		panic(err)
	}

	var interval, min, max, sym string

	sym = strings.ToLower(c.Param("symbol"))

	interval = c.Query("interval")
	if interval == "" {
		interval = controllers.INT_DAY
	}

	min = c.Query("min")
	if min == "" {
		min = "-inf"
	}

	max = c.Query("max")
	if max == "" {
		max = "+inf"
	}

	vals, err := db.ZRangeByScoreWithScores(controllers.CTX, sym, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
	}).Result()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Record not found!",
			"error":   err,
		})
		return
	}

	var sV, eV map[string]interface{}

	json.Unmarshal([]byte(vals[0].Member.(string)), &sV)
	json.Unmarshal([]byte(vals[len(vals)-1].Member.(string)), &eV)

	sDate := int(sV["timestamp"].(float64))
	eDate := int(eV["timestamp"].(float64))

	prices := []map[string]interface{}{}

	for timestamp := sDate; timestamp <= eDate; timestamp += controllers.Intervals[interval] {
		val, err := db.ZRangeByScoreWithScores(controllers.CTX, sym, &redis.ZRangeBy{
			Min:    fmt.Sprint(timestamp - (controllers.HOR_INTERVAL * 12)),
			Max:    fmt.Sprint(timestamp + (controllers.HOR_INTERVAL * 12)),
			Offset: 0,
		}).Result()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Record not found!",
				"error":   err,
			})
			return
		}

		if len(val) == 0 {
			continue
		}

		var member map[string]interface{}
		json.Unmarshal([]byte(val[0].Member.(string)), &member)

		member["timestampH"] = time.Unix(int64(timestamp), 0)
		member["timestampE"] = timestamp

		delete(member, "timestamp")

		prices = append(prices, member)
	}

	c.JSON(http.StatusOK, gin.H{
		"prices": prices,
	})
}
