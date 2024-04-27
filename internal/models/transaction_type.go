package models

import "gorm.io/gorm"

type TransactionStatusType string

const (
	Income  TransactionStatusType = "INCOME"
	Expense TransactionStatusType = "EXPENSE"
)

type TransactionType struct {
	gorm.Model
	Name TransactionStatusType
}
