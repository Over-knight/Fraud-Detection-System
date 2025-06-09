package model

// Transaction represents a financial transaction.
type Transaction struct {
	ID     string  `json:"id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}