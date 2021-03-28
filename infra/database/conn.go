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

	// Run seeds
	// seed(db)

	return db
}

// Seed puts seed data into the datbase
func seed(db *gorm.DB) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	indexPath := fmt.Sprintf("%s/assets/indexes", path)

	files, err := ioutil.ReadDir(indexPath)
	if err != nil {
		panic(err)
	}

	funds := []models.Fund{}

	for _, f := range files {
		networkFundPath := fmt.Sprintf("%s/%s/index.json", indexPath, f.Name())

		data, err := ioutil.ReadFile(networkFundPath)
		if err != nil {
			panic(err)
		}

		var body map[string]interface{}
		json.Unmarshal(data, &body)

		for _, i := range body["Indexes"].([]interface{}) {
			idx := i.(map[string]interface{})

			assets := []models.Asset{}

			for _, a := range idx["Assets"].([]interface{}) {
				ast := a.(map[string]interface{})

				assets = append(assets, models.Asset{
					Name:           ast["name"].(string),
					Symbol:         ast["symbol"].(string),
					Address:        ast["address"].(string),
					Decimals:       ast["decimals"].(float64),
					AllocationGwei: ast["initialAllocationGwei"].(string),
				})
			}

			funds = append(funds, models.Fund{
				Name:    idx["Name"].(string),
				Slug:    idx["Slug"].(string),
				Address: idx["Address"].(string),
				Network: body["Network"].(string),
				Assets:  assets,
			})
		}
	}

	for i := 1; i < len(funds); i++ {
		if err := db.Debug().Model(&models.Fund{}).Create(&funds[i]).Error; err != nil {
			log.Fatalf("cannot seed funds table: %v", err)
		}
	}
}
