package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Transfer struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FromAccount    Account      `json:"from_account"`
	FromAccountID  uint32    `gorm:"not null;" json:"from_account_id"`
	ToAccount    Account      `json:"to_account"`
	ToAccountID  uint32    `gorm:"not null;" json:"to_account_id"`
	Amount  uint32    `gorm:"not null;" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *Transfer) Prepare() {
	t.ID = 0
	t.FromAccount = Account{}
	t.ToAccount = Account{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Transfer) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if (t.FromAccountID == 0) {
		return errors.New("destiny account is invalid")
	}

	if (t.ToAccountID == 0) {
		return errors.New("destiny account is invalid")
	}

	return nil
}

func (t *Transfer) Save(db *gorm.DB) (*Transfer, error) {
	var err error

	err = db.Debug().Model(&Transfer{}).Create(&t).Error
	if err != nil {
		return &Transfer{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&Account{}).Where("id = ?", t.FromAccountID).Take(&t.FromAccount).Error
		if err != nil {
			return &Transfer{}, err
		}
	}

	return t, nil
}

func (t *Transfer) FindAll(db *gorm.DB) (*[]Transfer, error) {
	var err error

	transfers := []Transfer{}
	err = db.Debug().Model(&Transfer{}).Limit(100).Find(&transfers).Error
	if err != nil {
		return &[]Transfer{}, err
	}
	if len(transfers) > 0 {
		for i := range transfers {
			err = db.Debug().Model(&Account{}).Where("id = ?", transfers[i].FromAccountID).Take(&transfers[i].FromAccount).Error
			if err != nil {
				return &[]Transfer{}, err
			}

			err = db.Debug().Model(&Account{}).Where("id = ?", transfers[i].ToAccountID).Take(&transfers[i].ToAccount).Error
			if err != nil {
				return &[]Transfer{}, err
			}
		}
	}

	return &transfers, nil
}

func (t *Transfer) FindByOwnerId(db *gorm.DB, uid uint32) (*[]Transfer, error) {
	var err error

	transfers := []Transfer{}
	err = db.Debug().Model(&Transfer{}).Limit(100).Find(&transfers).Error
	if err != nil {
		return &[]Transfer{}, err
	}
	if len(transfers) > 0 {
		for i := range transfers {
			err := db.Debug().Model(&Account{}).Where("id = ?", transfers[i].FromAccountID).Take(&transfers[i].FromAccount).Error
			if err != nil {
				return &[]Transfer{}, err
			}
		}
	}

	return &transfers, nil
}

func (t *Transfer) FindByID(db *gorm.DB, pid uint64) (*Transfer, error) {
	var err error
	err = db.Debug().Model(&Transfer{}).Where("id = ?", pid).Take(&t).Error
	if err != nil {
		return &Transfer{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&Account{}).Where("id = ?", t.FromAccountID).Take(&t.FromAccount).Error
		if err != nil {
			return &Transfer{}, err
		}
	}
	return t, nil
}

func (t *Transfer) Update(db *gorm.DB) (*Transfer, error) {

	var err error

	err = db.Debug().Model(&Transfer{}).Where("id = ?", t.ID).Updates(Transfer{Amount: t.Amount, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Transfer{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&Account{}).Where("id = ?", t.FromAccountID).Take(&t.FromAccount).Error
		if err != nil {
			return &Transfer{}, err
		}
	}
	return t, nil
}

func (t *Transfer) Delete(db *gorm.DB, tid uint64, aid uint32) (int64, error) {

	db = db.Debug().Model(&Transfer{}).Where("id = ? and from_account_id = ?", tid, aid).Take(&Transfer{}).Delete(&Transfer{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Transfer not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}


