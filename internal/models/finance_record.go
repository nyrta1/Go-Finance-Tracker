package models

import "gorm.io/gorm"

type FinanceRecord struct {
	gorm.Model
	UserID            uint `gorm:"index"`
	Amount            float64
	TransactionTypeID uint
	TransactionType   TransactionStatusType
	CategoryID        uint
	Category          Category
	Note              string
}
