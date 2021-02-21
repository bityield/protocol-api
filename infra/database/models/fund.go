package models

import (
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Fund schema
type Fund struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"nickname"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare sets default attributes on model
func (f *Fund) Prepare() {
	f.ID = 0
	f.Name = html.EscapeString(strings.TrimSpace(f.Name))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

// SaveFund saves the model into the datbase
func (f *Fund) SaveFund(db *gorm.DB) (*Fund, error) {

	var err error
	err = db.Debug().Create(&f).Error
	if err != nil {
		return &Fund{}, err
	}
	return f, nil
}

// TableName sets the required table name in the database
func (Fund) TableName() string {
	return "funds"
}
