package dto

type FinanceRecordInput struct {
	Amount            float64 `json:"amount"`
	TransactionTypeID uint    `json:"transactionTypeID"`
	CategoryID        uint    `json:"categoryID"`
	Note              string  `json:"note"`
}
