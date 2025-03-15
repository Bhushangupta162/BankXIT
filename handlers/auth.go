// handlers/auth.go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "github.com/bhushangupta162/bank_management/models"
    "github.com/bhushangupta162/bank_management/utils"
)

// SignUpHandler handles user registration.
func SignUpHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Hash the password.
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
            return
        }
        user.Password = string(hashedPassword)

        // Save the user to the database.
        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
    }
}

// LoginHandler handles user login.
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Email    string `json:"email" binding:"required"`
            Password string `json:"password" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var user models.User
        if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
            return
        }

        // Generate a JWT token.
        tokenString, err := utils.GenerateToken(user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": tokenString})
    }
}
