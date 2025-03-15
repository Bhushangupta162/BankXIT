// models/transaction.go
package models

import (
    "time"

    "gorm.io/gorm"
)

type Transaction struct {
    ID              uint           `gorm:"primaryKey" json:"id"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

    AccountID       uint    `json:"account_id"`          // Which account this transaction is for
    TransactionType string  `json:"transaction_type"`    // "deposit", "withdrawal", "transfer"
    Amount          float64 `json:"amount"`              // How much money was moved
    Description     string  `json:"description"`         // Optional notes or reason
}
