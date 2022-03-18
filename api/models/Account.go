package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name  string    `gorm:"size:255;not null;" json:"name"`
	Description  string    `gorm:"size:255;" json:"description"`
	Owner    User      `json:"owner"`
	OwnerID  uint32    `gorm:"not null" json:"owner_id"`
	Balance  uint32    `gorm:"default:0;" json:"balance"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (a *Account) Prepare() {
	a.ID = 0
	a.Name = ""
	a.Description = ""
	a.Owner = User{}
	a.Balance = 0
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *Account) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (a *Account) Save(db *gorm.DB) (*Account, error) {
	var err error
	err = db.Debug().Model(&Account{}).Create(&a).Error
	if err != nil {
		return &Account{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.OwnerID).Take(&a.Owner).Error
		if err != nil {
			return &Account{}, err
		}
	}
	return a, nil
}

func (a *Account) FindAll(db *gorm.DB) (*[]Account, error) {
	var err error
	accounts := []Account{}
	err = db.Debug().Model(&Account{}).Limit(100).Find(&accounts).Error
	if err != nil {
		return &[]Account{}, err
	}
	if len(accounts) > 0 {
		for i := range accounts {
			err := db.Debug().Model(&User{}).Where("id = ?", accounts[i].OwnerID).Take(&accounts[i].Owner).Error
			if err != nil {
				return &[]Account{}, err
			}
		}
	}
	return &accounts, nil
}

func (a *Account) FindByID(db *gorm.DB, pid uint64) (*Account, error) {
	var err error
	err = db.Debug().Model(&Account{}).Where("id = ?", pid).Take(&a).Error
	if err != nil {
		return &Account{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.OwnerID).Take(&a.Owner).Error
		if err != nil {
			return &Account{}, err
		}
	}
	return a, nil
}

func (a *Account) Update(db *gorm.DB) (*Account, error) {

	var err error

	err = db.Debug().Model(&Account{}).Where("id = ?", a.ID).Updates(Account{Balance: a.Balance, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Account{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.OwnerID).Take(&a.Owner).Error
		if err != nil {
			return &Account{}, err
		}
	}
	return a, nil
}

func (a *Account) Delete(db *gorm.DB, aid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Account{}).Where("id = ? and owner_id = ?", aid, uid).Take(&Account{}).Delete(&Account{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Account not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}