package models

import (
	"time"
)

type Entry struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Account    Account      `json:"account"`
	AccountID  uint32    `gorm:"not null;" json:"account_id"`
	Amount  uint32    `gorm:"not null;" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
