package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bityield/protocol-api/infra/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// var ctx = context.Background()

// Connect kind of explanatory
func Connect() *gorm.DB {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, pass)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// Fund
	if err := db.Debug().DropTableIfExists(&models.Fund{}).Error; err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	if err = db.Debug().AutoMigrate(&models.Fund{}).Error; err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// Asset
	if err := db.Debug().DropTableIfExists(&models.Asset{}).Error; err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	if err = db.Debug().AutoMigrate(&models.Asset{}).Error; err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// Price
	if err := db.Debug().DropTableIfExists(&models.Price{}).Error; err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	if err = db.Debug().AutoMigrate(&models.Price{}).Error; err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// Address
	if err := db.Debug().DropTableIfExists(&models.Address{}).Error; err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	if err = db.Debug().AutoMigrate(&models.Address{}).Error; err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// Run seeds
	seed(db)

	return db
}

// Seed puts seed data into the datbase
func seed(db *gorm.DB) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/assets/indexes/ropsten/index.json", path))
	if err != nil {
		panic(err)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		panic(err)
	}

	funds := []models.Fund{}

	for x, i := range body["Indexes"].([]interface{}) {
		idx := i.(map[string]interface{})

		assets := []models.Asset{}

		for _, a := range idx["Assets"].([]interface{}) {
			ast := a.(map[string]interface{})

			assets = append(assets, models.Asset{
				Name:               ast["name"].(string),
				Symbol:             ast["symbol"].(string),
				Address:            ast["address"].(string),
				Decimals:           ast["decimals"].(float64),
				AllocationGwei:     ast["initialAllocationGwei"].(string),
				AllocationMantissa: ast["initialAllocationMantissa"].(string),
			})
		}

		funds = append(funds, models.Fund{
			Name:    idx["Name"].(string),
			Slug:    idx["Slug"].(string),
			Icon:    idx["Icon"].(string),
			Address: idx["Address"].(string),
			Network: body["Network"].(string),
			Assets:  assets,
		})

		if err := db.Debug().Model(&models.Fund{}).Create(&funds[x]).Error; err != nil {
			log.Fatalf("cannot seed funds table: %v", err)
		}
	}
}
