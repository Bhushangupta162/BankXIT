// utils/jwt.go
package utils

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
)

// JwtSecret is our secret key for signing tokens.
// In a production app, store this securely (e.g., environment variable).
var JwtSecret = []byte("your_secret_key")

// GenerateToken generates a JWT token for a given user ID.
func GenerateToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(72 * time.Hour).Unix(), // Token valid for 72 hours
    })
    return token.SignedString(JwtSecret)
}
