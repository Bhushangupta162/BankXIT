// main.go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/bhushangupta162/bank_management/models"
    "github.com/bhushangupta162/bank_management/handlers"
    "github.com/bhushangupta162/bank_management/utils"
)

func main() {
    // Create a new Gin router.
    router := gin.Default()

    // Connect to PostgreSQL.
    // Make sure this DSN matches your PostgreSQL settings (using Docker or local installation).
    dsn := "host=postgres user=postgres password=postgres dbname=bank port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Migrate the User model.
    db.AutoMigrate(&models.User{})
    db.AutoMigrate(&models.User{}, &models.Account{})
    db.AutoMigrate(&models.User{},&models.Account{},&models.Transaction{},)
    db.AutoMigrate(&models.User{},&models.Account{},&models.Transaction{},&models.Loan{},)
    

    // Define the authentication routes.
    router.POST("/signup", handlers.SignUpHandler(db))
    router.POST("/login", handlers.LoginHandler(db))

    // Account routes
    router.POST("/accounts", handlers.CreateAccountHandler(db))             // Create an account
    router.GET("/accounts/:id", handlers.GetAccountHandler(db))
    router.POST("/accounts/:id/deposit", handlers.DepositHandler(db))
    router.POST("/accounts/:id/withdraw", handlers.WithdrawHandler(db))
    router.POST("/accounts/transfer", handlers.TransferHandler(db))
    router.GET("/accounts/:id/transactions", handlers.GetTransactionsHandler(db))

    // Loan endpoints
    router.POST("/loans/apply", handlers.ApplyLoanHandler(db))
    router.PATCH("/loans/:id/status", handlers.UpdateLoanStatusHandler(db))   // Approve/Reject
    router.POST("/loans/:id/repay", handlers.RepayLoanHandler(db))

    // Example protected route.
    router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "You are authorized!"})
    })

    // Start the server on port 8080.
    router.Run(":8080")
}

// AuthMiddleware verifies JWT tokens for protected endpoints.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
            return
        }

        // Parse and validate the token.
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return utils.JwtSecret, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Token is valid, continue.
        c.Next()
    }
}
