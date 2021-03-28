package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bityield/protocol-api/backend"
	"github.com/bityield/protocol-api/infra/database/models"
)

const (
	BASE_ENDPOINT = "https://web-api.coinmarketcap.com/v1/cryptocurrency/ohlcv/historical"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Execute(b *backend.Backend) error {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// The main source of truth for this is our coins.json file in root of this repository
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/coins.json", path))
	if err != nil {
		panic(err)
	}

	// beginTime should be a long time ago, endTime is today
	bTime, eTime := "1420099200", fmt.Sprint(time.Now().Unix())

	var coins map[string]interface{}
	json.Unmarshal(data, &coins)

	for _, val := range coins["coins"].([]interface{}) {
		obj := val.(map[string]interface{})

		nme := obj["name"].(string)
		sym := obj["symbol"].(string)
		url := fmt.Sprintf("%s?convert=USD&symbol=%s&time_start=%s&time_end=%s", BASE_ENDPOINT, sym, bTime, eTime)

		b.L.Infof("CoinMarketCap saving [%s] from: [%s]\n", sym, url)

		// Query endpoint
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		// Parse data
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var body map[string]interface{}
		json.Unmarshal(data, &body)

		for _, q := range body["data"].(map[string]interface{})["quotes"].([]interface{}) {
			tO := q.(map[string]interface{})["time_open"]
			tC := q.(map[string]interface{})["time_close"]

			s := q.(map[string]interface{})["quote"]
			u := s.(map[string]interface{})["USD"].(map[string]interface{})

			price := models.Price{
				Name:      nme,
				Symbol:    sym,
				TimeOpen:  tO.(string),
				TimeClose: tC.(string),
				Open:      u["open"].(float64),
				High:      u["high"].(float64),
				Low:       u["low"].(float64),
				Close:     u["close"].(float64),
				MarketCap: u["market_cap"].(float64),
				Volume:    u["volume"].(float64),
			}

			if err := b.R.D.Debug().Model(&models.Price{}).Create(&price).Error; err != nil {
				panic(err)
			}
		}

		// Sleep for some random time, don't overload CMC
		duration := rand.Intn(15)
		b.L.Infof("Sleeping %d seconds...\n", duration)
		time.Sleep(time.Duration(duration) * time.Second)
	}

	return nil
}
