package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Address ...
type Address struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// Prepare ...
func (u *Address) Prepare() {
	u.ID = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// SaveAddress ...
func (u *Address) SaveAddress(db *gorm.DB) (*Address, error) {

	err := db.Debug().Create(&u).Error
	if err != nil {
		return &Address{}, err
	}
	return u, nil
}

// FindAllAddresses ...
func (u *Address) FindAllAddresses(db *gorm.DB) (*[]Address, error) {
	var err error
	addresses := []Address{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&addresses).Error
	if err != nil {
		return &[]Address{}, err
	}
	return &addresses, err
}

// FindAddressByID ...
func (u *Address) FindAddressByID(db *gorm.DB, uid uint32) (*Address, error) {
	err := db.Debug().Model(Address{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Address{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Address{}, errors.New("address not found")
	}
	return u, err
}

// UpdateAddress ...
func (u *Address) UpdateAddress(db *gorm.DB, uid uint32) (*Address, error) {
	var err error
	db = db.Debug().Model(&Address{}).Where("id = ?", uid).Take(&Address{}).UpdateColumns(
		map[string]interface{}{
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Address{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Address{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Address{}, err
	}
	return u, nil
}

// DeleteAddress ...
func (u *Address) DeleteAddress(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Address{}).Where("id = ?", uid).Take(&Address{}).Delete(&Address{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
