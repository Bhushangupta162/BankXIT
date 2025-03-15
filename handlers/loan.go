// handlers/loan.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/bhushangupta162/bank_management/models"
)

// ApplyLoanHandler - a user requests a new loan
func ApplyLoanHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            UserID       uint    `json:"user_id" binding:"required"`
            Principal    float64 `json:"principal" binding:"required"`
            InterestRate float64 `json:"interest_rate" binding:"required"`
            TermMonths   int     `json:"term_months" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        loan := models.Loan{
            UserID:            input.UserID,
            Principal:         input.Principal,
            InterestRate:      input.InterestRate,
            TermMonths:        input.TermMonths,
            Status:            "pending",
            OutstandingBalance: input.Principal,
        }

        if err := db.Create(&loan).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create loan"})
            return
        }

        c.JSON(http.StatusCreated, loan)
    }
}

// UpdateLoanStatusHandler - admin approves or rejects a loan
func UpdateLoanStatusHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        loanIDStr := c.Param("id")
        loanID, err := strconv.Atoi(loanIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
            return
        }

        var input struct {
            Status string `json:"status" binding:"required"` // e.g. "approved", "rejected"
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var loan models.Loan
        if err := db.First(&loan, loanID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
            return
        }

        // Only allow transitions from pending -> [approved or rejected], for example
        if loan.Status != "pending" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Loan is not pending"})
            return
        }

        if input.Status == "approved" {
            loan.Status = "active" // or "approved"
            // Optionally set a start date or schedule interest calculation
        } else if input.Status == "rejected" {
            loan.Status = "rejected"
            // No further changes needed
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
            return
        }

        if err := db.Save(&loan).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update loan status"})
            return
        }

        c.JSON(http.StatusOK, loan)
    }
}

// RepayLoanHandler - user repays part of a loan
func RepayLoanHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        loanIDStr := c.Param("id")
        loanID, err := strconv.Atoi(loanIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
            return
        }

        var input struct {
            Amount float64 `json:"amount" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var loan models.Loan
        if err := db.First(&loan, loanID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
            return
        }

        // Check status
        if loan.Status != "active" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Loan is not active for repayment"})
            return
        }

        if input.Amount <= 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Repayment amount must be positive"})
            return
        }

        if input.Amount > loan.OutstandingBalance {
            input.Amount = loan.OutstandingBalance // or reject if you want strict
        }

        // Deduct from the outstanding balance
        loan.OutstandingBalance -= input.Amount
        if loan.OutstandingBalance <= 0 {
            loan.OutstandingBalance = 0
            loan.Status = "closed" // fully repaid
        }

        if err := db.Save(&loan).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan"})
            return
        }

        // Optionally, log a transaction in your transactions table
        // e.g., "loan repayment" if you want a record
        // We can do that if you want a separate model referencing "loans" or do it in "transactions" with a new field.
        // For simplicity, let's skip or do something simple:
        // (Assuming you have a Transaction model, you can reuse that or create a new type)

        c.JSON(http.StatusOK, loan)
    }
}



