package models

import (
	"time"
)

type Account struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Owner    User      `json:"owner"`
	OwnerID  uint32    `gorm:"not null" json:"owner_id"`
	Balance  uint32    `gorm:"default:0;" json:"balance"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (a *Account) Prepare() {
	a.ID = 0
	a.Owner = User{}
	a.Balance = 0
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}
