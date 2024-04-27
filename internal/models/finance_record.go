package models

import "gorm.io/gorm"

type FinanceRecord struct {
	gorm.Model
	UserID            uint            `gorm:"index" json:"userID"`
	Amount            float64         `json:"amount"`
	TransactionTypeID uint            `json:"transactionTypeID"`
	TransactionType   TransactionType `gorm:"foreignKey:TransactionTypeID"`
	CategoryID        uint            `json:"categoryID"`
	Category          Category        `gorm:"foreignKey:CategoryID"`
	Note              string          `json:"note"`
}
