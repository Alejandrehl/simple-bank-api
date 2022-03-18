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
	Amount  uint32    `gorm:"default:0;not null;" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (e *Entry) Prepare() {
	e.ID = 0
	e.Account = Account{}
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
		err = db.Debug().Model(&Account{}).Where("id = ?", e.AccountID).Take(&e.Account).Error
		if err != nil {
			return &Entry{}, err
		}
	}

	return e, nil
}

func (e *Entry) FindAll(db *gorm.DB, uid uint32) (*[]Entry, error) {
	var err error

	entries := []Entry{}
	err = db.Debug().Model(&Entry{}).Limit(100).Find(&entries).Error
	if err != nil {
		return &[]Entry{}, err
	}
	if len(entries) > 0 {
		for i := range entries {
			err := db.Debug().Model(&Account{}).Where("id = ?", entries[i].AccountID).Take(&entries[i].Account).Error
			if err != nil {
				return &[]Entry{}, err
			}
		}
	}

	return &entries, nil
}

func (e *Entry) FindByID(db *gorm.DB, eid uint64, uid uint32) (*Entry, error) {
	var err error

	// err = db.Debug().Model(&Entry{}).Where("id = ?", eid).Take(&e).Error
	err = db.Debug().Model(&Entry{}).Joins("Account", db.Where(&Account{OwnerID: uid})).Where("id = ?", eid).Take(&e).Error
	if err != nil {
		return &Entry{}, err
	}
	if e.ID != 0 {
		err = db.Debug().Model(&Account{}).Where("id = ?", e.AccountID).Take(&e.Account).Error
		if err != nil {
			return &Entry{}, err
		}
	}

	return e, nil
}

func (e *Entry) Update(db *gorm.DB) (*Entry, error) {
	var err error

	err = db.Debug().Model(&Entry{}).Where("id = ?", e.ID).Updates(Entry{Amount: e.Amount, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Entry{}, err
	}
	if e.ID != 0 {
		err = db.Debug().Model(&Account{}).Where("id = ?", e.AccountID).Take(&e.Account).Error
		if err != nil {
			return &Entry{}, err
		}
	}

	return e, nil
}

func (e *Entry) Delete(db *gorm.DB, eid uint64, aid uint32) (int64, error) {
	db = db.Debug().Model(&Entry{}).Where("id = ? and account_id = ?", eid, aid).Take(&Entry{}).Delete(&Entry{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Entry not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
