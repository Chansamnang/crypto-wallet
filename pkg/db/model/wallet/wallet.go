package wallet

import (
	"gorm.io/gorm"
)

func NewWallet(db *gorm.DB) *Wallet {
	return &Wallet{
		DB: db,
	}
}

type Wallet struct {
	DB *gorm.DB
}
