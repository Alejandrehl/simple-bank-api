package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Entry struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Account    Account      `json:"account"`
	AccountID  uint32    `gorm:"not null;" json:"account_id"`
	Amount  uint32    `gorm:"not null;" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (e *Entry) Prepare() {
	e.ID = 0
	e.Owner = User{}
	e.Balance = 0
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

func (e *Entry) Validate() error {
	return nil
}

func (e *Entry) Save(db *gorm.DB) (*Entry, error) {
	var err error
	err = db.Debug().Model(&Entry{}).Create(&e).Error
	if err != nil {
		return &Entry{}, err
	}
	if e.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", e.OwnerID).Take(&e.Owner).Error
		if err != nil {
			return &Entry{}, err
		}
	}
	return e, nil
}

func (e *Entry) FindAll(db *gorm.DB) (*[]Entry, error) {
	var err error
	Entrys := []Entry{}
	err = db.Debug().Model(&Entry{}).Limit(100).Find(&Entrys).Error
	if err != nil {
		return &[]Entry{}, err
	}
	if len(Entrys) > 0 {
		for i := range Entrys {
			err := db.Debug().Model(&User{}).Where("id = ?", Entrys[i].OwnerID).Take(&Entrys[i].Owner).Error
			if err != nil {
				return &[]Entry{}, err
			}
		}
	}
	return &Entrys, nil
}

func (e *Entry) FindByID(db *gorm.DB, pid uint64) (*Entry, error) {
	var err error
	err = db.Debug().Model(&Entry{}).Where("id = ?", pid).Take(&e).Error
	if err != nil {
		return &Entry{}, err
	}
	if e.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", e.OwnerID).Take(&e.Owner).Error
		if err != nil {
			return &Entry{}, err
		}
	}
	return e, nil
}

func (e *Entry) Update(db *gorm.DB) (*Entry, error) {

	var err error

	err = db.Debug().Model(&Entry{}).Where("id = ?", e.ID).Updates(Entry{Balance: e.Balance, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Entry{}, err
	}
	if e.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", e.OwnerID).Take(&e.Owner).Error
		if err != nil {
			return &Entry{}, err
		}
	}
	return e, nil
}

func (e *Entry) Delete(db *gorm.DB, aid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Entry{}).Where("id = ? and owner_id = ?", aid, uid).Take(&Entry{}).Delete(&Entry{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Entry not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
