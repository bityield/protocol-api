package models

import (
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Asset schema
type Price struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;" json:"name"`
	Symbol    string    `gorm:"size:255;not null;" json:"symbol"`
	TimeOpen  string    `gorm:"size:255;not null;" json:"timeOpen"`
	TimeClose string    `gorm:"size:255;not null;" json:"timeClose"`
	Open      float64   `gorm:"type:decimal(18,10);" json:"open"`
	High      float64   `gorm:"type:decimal(18,10);" json:"high"`
	Low       float64   `gorm:"type:decimal(18,10);" json:"low"`
	Close     float64   `gorm:"type:decimal(18,10);" json:"close"`
	Volume    float64   `gorm:"type:decimal(18,2);" json:"volume"`
	MarketCap float64   `gorm:"type:decimal(18,2);" json:"marketCap"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// Prepare sets default attributes on model
func (f *Price) Prepare() {
	f.ID = 0
	f.Name = html.EscapeString(strings.TrimSpace(f.Name))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

// SaveAsset saves the model into the datbase
func (f *Price) SaveAsset(db *gorm.DB) (*Price, error) {
	var err error

	if err = db.Debug().Create(&f).Error; err != nil {
		return &Price{}, err
	}

	return f, nil
}

// TableName sets the required table name in the database
func (Price) TableName() string {
	return "prices"
}
