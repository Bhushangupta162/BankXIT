// models/account.go
package models

import (
    "time"

    "gorm.io/gorm"
)

type Account struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    UserID  uint    `json:"user_id"`                      // Foreign key to User
    Balance float64 `json:"balance" gorm:"not null;default:0"`
    // Add more fields like AccountType if needed
}
