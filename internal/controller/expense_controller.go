package controller

import (
	"expense-mgmt/db"
	"expense-mgmt/internal/models"
	"expense-mgmt/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateExpense creates a new expense record
func CreateExpense(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Get the user ID from the JWT token in the middleware (assumed to be set)
	userID, exists := c.Get("userId")
	if !exists {
		utils.SendResponse(c, http.StatusUnauthorized, "User not authorized", nil, nil)
		return
	}

	var expense models.Expense
	// Bind the incoming JSON payload to the Expense model
	if err := c.ShouldBindJSON(&expense); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid input", nil, nil)
		return
	}

	// Validate input
	if expense.Amount <= 0 {
		utils.SendResponse(c, http.StatusBadRequest, "Amount must be positive", nil, nil)
		return
	}

	if expense.Date.After(time.Now()) {
		utils.SendResponse(c, http.StatusBadRequest, "Date cannot be in the future", nil, nil)
		return
	}

	// Set the user_id
	expense.UserID = userID.(uuid.UUID)

	// Fetch category by category_id to ensure it exists
	var category models.Category
	if err := DB.First(&category, "id = ?", expense.CategoryID).Error; err != nil {
    utils.SendResponse(c, http.StatusBadRequest, "Invalid category ID", nil, nil)
    return
	}

	// Save the expense to the database
	if err := DB.Create(&expense).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Error saving expense", nil, nil)
		return
	}

	// Respond with the created expense
	utils.SendResponse(c, http.StatusOK, "Expense created successfully", expense, nil)
}


// ListUserExpenses retrieves all expenses for the authenticated user
func ListUserExpenses(c *gin.Context) {
	// Get the user_id from the context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Initialize filters and pagination
	var expenses []models.Expense
	var totalCount int64
	page := utils.ParseQueryInt(c, "page", 1)
	limit := utils.ParseQueryInt(c, "limit", 10)
	offset := (page - 1) * limit

	// Optional filters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	categoryID := c.Query("category_id")
	minAmount := c.Query("min_amount")
	maxAmount := c.Query("max_amount")
	sort := c.DefaultQuery("sort", "date")
	order := c.DefaultQuery("order", "asc")

	// Build the query with filters
	query := db.GetDBInstance().Model(&models.Expense{}).Where("user_id = ?", userID)

	if startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("date >= ?", parsedDate)
		}
	}
	if endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("date <= ?", parsedDate)
		}
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if minAmount != "" {
		query = query.Where("amount >= ?", minAmount)
	}
	if maxAmount != "" {
		query = query.Where("amount <= ?", maxAmount)
	}

	// Count the total number of records
	query.Count(&totalCount)

	// Apply sorting, pagination, and fetch the results
	if err := query.Order(sort + " " + order).
		Offset(offset).
		Limit(limit).
		Find(&expenses).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch expenses", nil, nil)
		return
	}

	// Use CalculatePagination to generate pagination info
	pagination := utils.CalculatePagination(int(totalCount), page, limit)

	// Send the response with pagination
	utils.SendResponse(c, http.StatusOK, "Expenses fetched successfully", expenses, pagination)
}

// GetExpense retrieves detailed information about a specific expense
func GetExpense(c *gin.Context) {
	// Get the user_id from the context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve the expenseId from the URL
	expenseID := c.Param("expenseId")
	fmt.Println(expenseID)

	// Initialize the expense variable
	var expense models.Expense

	// Fetch the expense by its ID for the authenticated user
	if err := db.GetDBInstance().Model(&models.Expense{}).
		Where("user_id = ? AND expense_id = ?", userID, expenseID).
		First(&expense).Error; err != nil {
		// Check if the expense was not found
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Expense not found", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch expense", nil, nil)
		}
		return
	}

	// Send the response with the expense details
	utils.SendResponse(c, http.StatusOK, "Expense fetched successfully", expense, nil)
}

// UpdateExpense updates an existing expense
func UpdateExpense(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found in request context", nil, nil)
		return
	}

	// Retrieve the expenseId from the URL
	expenseID := c.Param("expenseId")

	// Initialize the expense variable
	var expense models.Expense

	// Fetch the expense by its ID and user
	if err := db.GetDBInstance().Where("user_id = ? AND expense_id = ?", userID, expenseID).First(&expense).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Expense not found for the given user", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Database error while fetching expense", nil, nil)
		}
		return
	}

	// Bind the JSON request data
	var updateData models.Expense
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid JSON input: "+err.Error(), nil, nil)
		return
	}

	// Track fields to update and validate inputs
	updateFields := map[string]interface{}{}
	if updateData.Amount != 0 {
		if updateData.Amount < 0 {
			utils.SendResponse(c, http.StatusBadRequest, "Amount must be a positive value", nil, nil)
			return
		}
		updateFields["amount"] = updateData.Amount
	}

	if !updateData.Date.IsZero() {
		updateFields["date"] = updateData.Date
	}

	if updateData.Description != "" {
		updateFields["description"] = updateData.Description
	}

	if updateData.CategoryID != uuid.Nil {
		// Check if the provided CategoryID exists
		var category models.Category
		if err := db.GetDBInstance().Where("id = ?", updateData.CategoryID).First(&category).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.SendResponse(c, http.StatusBadRequest, "Invalid category ID", nil, nil)
				return
			} else {
				utils.SendResponse(c, http.StatusInternalServerError, "Database error while validating category", nil, nil)
				return
			}
		}
		updateFields["category_id"] = updateData.CategoryID
	}

	if updateData.ReceiptID != nil {
		// Check if the provided ReceiptID exists
		var receipt models.Receipt
		if err := db.GetDBInstance().Where("receipt_id = ?", updateData.ReceiptID).First(&receipt).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.SendResponse(c, http.StatusBadRequest, "Invalid receipt ID", nil, nil)
				return
			} else {
				utils.SendResponse(c, http.StatusInternalServerError, "Database error while validating receipt", nil, nil)
				return
			}
		}
		updateFields["receipt_id"] = updateData.ReceiptID
	}

	// Apply the updates
	if len(updateFields) > 0 {
		if err := db.GetDBInstance().Model(&expense).Updates(updateFields).Error; err != nil {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to update expense in database", nil, nil)
			return
		}
	} else {
		utils.SendResponse(c, http.StatusBadRequest, "No valid fields to update", nil, nil)
		return
	}

	// Fetch the updated expense to send a fresh response
	if err := db.GetDBInstance().First(&expense, "expense_id = ?", expenseID).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch updated expense", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Expense updated successfully", expense, nil)
}

// DeleteExpense deletes an expense and its associated receipt if it exists
func DeleteExpense(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve the expenseId from the URL
	expenseID := c.Param("expenseId")

	// Initialize the expense variable
	var expense models.Expense

	// Fetch the expense by its ID and user
	if err := db.GetDBInstance().Where("user_id = ? AND expense_id = ?", userID, expenseID).First(&expense).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Expense not found", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch expense", nil, nil)
		}
		return
	}

	// Check if there’s an associated receipt and delete it if it exists
	if expense.ReceiptID != nil {
		if err := db.GetDBInstance().Where("receipt_id = ?", expense.ReceiptID).Delete(&models.Receipt{}).Error; err != nil {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete associated receipt", nil, nil)
			return
		}
	}

	// Delete the expense
	if err := db.GetDBInstance().Delete(&expense).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete expense", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Expense deleted successfully", nil, nil)
}

// ExpenseAnalysis calculates and returns insights about user expenses
func ExpenseAnalysis(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve query parameters
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	categoryID := c.DefaultQuery("category_id", "")
	period := c.DefaultQuery("period", "month")

	// Build the query
	var analysisData []models.Expense
	query := db.GetDBInstance().Where("user_id = ?", userID)

	// Apply filters
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Group by period (e.g., monthly)
	if period == "month" {
		query = query.Group("strftime('%Y-%m', date)").Select("strftime('%Y-%m', date) as period, sum(amount) as total_spending").Scan(&analysisData)
	}

	// Execute query
	if err := query.Find(&analysisData).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch analysis data", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Expense analysis fetched successfully", analysisData, nil)
}
