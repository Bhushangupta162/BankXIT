// handlers/account.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/bhushangupta162/bank_management/models"
)

// CreateAccountHandler creates a new account for a specific user
func CreateAccountHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // We expect a JSON body with { "user_id": <number> }
        var input struct {
            UserID uint `json:"user_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Initialize an Account model
        account := models.Account{
            UserID:  input.UserID,
            Balance: 0, // default balance
        }

        if err := db.Create(&account).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, account)
    }
}

// GetAccountHandler retrieves an account by ID
func GetAccountHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // account ID from the URL param
        accountIDStr := c.Param("id")
        accountID, err := strconv.Atoi(accountIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
            return
        }

        var account models.Account
        if err := db.First(&account, accountID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
            return
        }

        c.JSON(http.StatusOK, account)
    }
}

// DepositHandler deposits a given amount into an account
func DepositHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Amount float64 `json:"amount" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        accountIDStr := c.Param("id")
        accountID, err := strconv.Atoi(accountIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
            return
        }

        var account models.Account
        if err := db.First(&account, accountID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
            return
        }

		if input.Amount <= 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
            return
        }

        // Perform deposit
        account.Balance += input.Amount
        if err := db.Save(&account).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

		// Log the transaction
        txRecord := models.Transaction{
            AccountID:       account.ID,
            TransactionType: "deposit",
            Amount:          input.Amount,
            Description:     "Deposit operation",
        }
        if err := db.Create(&txRecord).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transaction"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "account":     account,
            "transaction": txRecord,
        })
    }
}

// WithdrawHandler withdraws a given amount from an account and logs a transaction
func WithdrawHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Amount float64 `json:"amount" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        accountIDStr := c.Param("id")
        accountID, err := strconv.Atoi(accountIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
            return
        }

        var account models.Account
        if err := db.First(&account, accountID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
            return
        }

        if input.Amount <= 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
            return
        }

        // Check balance
        if account.Balance < input.Amount {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
            return
        }

        // Perform withdrawal
        account.Balance -= input.Amount
        if err := db.Save(&account).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Log the transaction
        txRecord := models.Transaction{
            AccountID:       account.ID,
            TransactionType: "withdrawal",
            Amount:          input.Amount,
            Description:     "Withdrawal operation",
        }
        if err := db.Create(&txRecord).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transaction"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "account":     account,
            "transaction": txRecord,
        })
    }
}

// TransferHandler transfers an amount from one account to another in a single transaction
func TransferHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            FromAccountID uint    `json:"from_account_id" binding:"required"`
            ToAccountID   uint    `json:"to_account_id" binding:"required"`
            Amount        float64 `json:"amount" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if input.FromAccountID == input.ToAccountID {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to the same account"})
            return
        }

        if input.Amount <= 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
            return
        }

        // Use a DB transaction to ensure both steps succeed or fail together
        tx := db.Begin()

        var fromAccount, toAccount models.Account
        // Lock rows "FOR UPDATE" if you want concurrency safety
        if err := tx.First(&fromAccount, input.FromAccountID).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "Source account not found"})
            return
        }

        if err := tx.First(&toAccount, input.ToAccountID).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "Destination account not found"})
            return
        }

        // Check balance
        if fromAccount.Balance < input.Amount {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
            return
        }

        // Perform transfer
        fromAccount.Balance -= input.Amount
        toAccount.Balance += input.Amount

        if err := tx.Save(&fromAccount).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if err := tx.Save(&toAccount).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Log transactions
        fromTx := models.Transaction{
            AccountID:       fromAccount.ID,
            TransactionType: "transfer-out",
            Amount:          input.Amount,
            Description:     "Transfer to account " + strconv.Itoa(int(toAccount.ID)),
        }
        if err := tx.Create(&fromTx).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transfer-out transaction"})
            return
        }
        toTx := models.Transaction{
            AccountID:       toAccount.ID,
            TransactionType: "transfer-in",
            Amount:          input.Amount,
            Description:     "Transfer from account " + strconv.Itoa(int(fromAccount.ID)),
        }
        if err := tx.Create(&toTx).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transfer-in transaction"})
            return
        }

        // Commit transaction
        if err := tx.Commit().Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "from_account": fromAccount,
            "to_account":   toAccount,
            "out_tx":       fromTx,
            "in_tx":        toTx,
        })
    }
}

// GetTransactionsHandler returns a list of transactions for an account
func GetTransactionsHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        accountIDStr := c.Param("id")
        accountID, err := strconv.Atoi(accountIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
            return
        }

        var transactions []models.Transaction
        if err := db.Where("account_id = ?", accountID).Order("created_at DESC").Find(&transactions).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch transactions"})
            return
        }

        c.JSON(http.StatusOK, transactions)
    }
}




