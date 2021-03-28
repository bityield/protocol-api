package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bityield/protocol-api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()

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

	intervals = map[string]int{
		INT_HOR: HOR_INTERVAL,
		INT_DAY: DAY_INTERVAL,
		INT_WEK: WEK_INTERVAL,
		INT_MNT: MNT_INTERVAL,
		INT_YER: YER_INTERVAL,
	}
)

// GetHistoricals ...
func GetHistoricals(c *gin.Context) {
	db, err := controllers.GetRedis(c)
	if err != nil {
		panic(err)
	}

	var interval, min, max, sym string

	sym = strings.ToLower(c.Param("symbol"))

	interval = c.Query("interval")
	if interval == "" {
		interval = INT_DAY
	}

	min = c.Query("min")
	if min == "" {
		min = "-inf"
	}

	max = c.Query("max")
	if max == "" {
		max = "+inf"
	}

	vals, err := db.ZRangeByScoreWithScores(ctx, sym, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
	}).Result()

	if err != nil {
		panic(err)
	}

	var sV, eV map[string]interface{}

	json.Unmarshal([]byte(vals[0].Member.(string)), &sV)
	json.Unmarshal([]byte(vals[len(vals)-1].Member.(string)), &eV)

	sDate := int(sV["timestamp"].(float64))
	eDate := int(eV["timestamp"].(float64))

	prices := []map[string]interface{}{}

	for timestamp := sDate; timestamp <= eDate; timestamp += intervals[interval] {
		// fmt.Println("Timestamp current:", timestamp, ", DateTime:", time.Unix(int64(timestamp), 0))

		val, err := db.ZRangeByScoreWithScores(ctx, sym, &redis.ZRangeBy{
			Min:    fmt.Sprint(timestamp - (HOR_INTERVAL * 12)),
			Max:    fmt.Sprint(timestamp + (HOR_INTERVAL * 12)),
			Offset: 0,
		}).Result()

		if err != nil {
			panic(err)
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
