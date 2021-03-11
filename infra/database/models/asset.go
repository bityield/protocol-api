package models

import (
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Asset schema
type Asset struct {
	ID             uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name           string    `gorm:"size:255;not null;" json:"name"`
	Symbol         string    `gorm:"size:255;not null;" json:"symbol"`
	Address        string    `gorm:"size:255;not null;" json:"address"`
	Decimals       float64   `json:"decimals"`
	AllocationGwei string    `json:"initialAllocationGwei"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Associations
	FundID uint
}

// Prepare sets default attributes on model
func (f *Asset) Prepare() {
	f.ID = 0
	f.Name = html.EscapeString(strings.TrimSpace(f.Name))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

// SaveAsset saves the model into the datbase
func (f *Asset) SaveAsset(db *gorm.DB) (*Asset, error) {

	var err error
	err = db.Debug().Create(&f).Error
	if err != nil {
		return &Asset{}, err
	}
	return f, nil
}

// TableName sets the required table name in the database
func (Asset) TableName() string {
	return "assets"
}
