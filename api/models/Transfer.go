package models

import (
	"time"
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
