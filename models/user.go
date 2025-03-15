// models/user.go
package models

import "gorm.io/gorm"

// User represents a user in the system.
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}
