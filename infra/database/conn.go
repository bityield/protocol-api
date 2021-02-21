package database

import (
	"fmt"
	"log"

	"github.com/bityield/bityield-api/infra/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ConnectDatabase kind of explanatory
func ConnectDatabase() *gorm.DB {
	dcs := "host=localhost port=5432 user=postgres dbname=bityield-api sslmode=disable password="
	db, err := gorm.Open("postgres", dcs)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to [mysql] database"))
	}

	if err := db.Debug().DropTableIfExists(&models.Fund{}).Error; err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	if err = db.Debug().AutoMigrate(&models.Fund{}).Error; err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// Run seeds
	seed(db)

	return db
}

// Seed puts seed data into the datbase
func seed(db *gorm.DB) {
	funds := []models.Fund{
		{
			Name: "GeneralPurposeV1",
		},
		{
			Name: "GeneralPurposeV2",
		},
	}

	for i, _ := range funds {
		if err := db.Debug().Model(&models.Fund{}).Create(&funds[i]).Error; err != nil {
			log.Fatalf("cannot seed funds table: %v", err)
		}

		// posts[i].AuthorID = users[i].ID

		// err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		// if err != nil {
		// 	log.Fatalf("cannot seed posts table: %v", err)
		// }
	}

}
