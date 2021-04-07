package funds

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bityield/protocol-api/controllers"
	"github.com/gin-gonic/gin"
)

// GetFundHistoricals ...
func GetFundHistoricals(c *gin.Context) {
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

	vals, err := db.Get(controllers.CTX, sym).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Record not found!",
			"error":   err,
		})
		return
	}

	var body map[string]interface{}

	json.Unmarshal([]byte(vals), &body)

	c.JSON(http.StatusOK, body)
}
