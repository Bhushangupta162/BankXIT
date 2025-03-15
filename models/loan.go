// models/loan.go
package models

import (
    "time"

    "gorm.io/gorm"
)

// Loan represents a credit/loan that a user can apply for.
type Loan struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

    UserID       uint    `json:"user_id"`            // Which user the loan belongs to
    Principal    float64 `json:"principal"`          // Original amount
    InterestRate float64 `json:"interest_rate"`      // Annual interest rate (e.g., 5.0 = 5%)
    TermMonths   int     `json:"term_months"`        // For example, 12, 24, 36 months, etc.
    Status       string  `json:"status"`             // e.g. "pending", "approved", "rejected", "active", "closed"

    // You might track extra fields:
    OutstandingBalance float64 `json:"outstanding_balance"` // How much is left to repay
    // You can also add fields like monthlyPayment, nextPaymentDue, etc. as needed
}
