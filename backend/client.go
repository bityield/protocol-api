package backend

import (

	// "strconv"

	"github.com/bityield/protocol-api/config"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"

	db "github.com/bityield/protocol-api/infra/database"
	rd "github.com/bityield/protocol-api/infra/redis"

	"github.com/sirupsen/logrus"
)

// Repository ...
type Repository struct {
	D *gorm.DB
	R *redis.Client
}

// Backend - main struct for the entire application configuration
type Backend struct {
	// C - contains the yaml file configuration key/values and other env specifics
	C config.Reader

	// L - a logrus logger, customized for this application
	L *logrus.Logger

	// R - a repository object for holding db/redis connections
	R *Repository
}

// NewBackend - factory method for producing a new type of Backend
func NewBackend() (*Backend, error) {
	c, err := config.LoadConfig(config.Defaults)
	if err != nil {
		return nil, err
	}

	// Database connection
	ds := db.Connect()

	// Redis connection
	rd := rd.Connect()

	// Base BackendConfiguration to link structs and objects
	var bc = &Backend{
		C: c,
		L: config.LoadLogger(c),
		R: &Repository{
			D: ds,
			R: rd,
		},
	}

	return bc, nil
}
