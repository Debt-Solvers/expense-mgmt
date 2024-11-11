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

// CreateBudget creates a new budget for the user
func CreateBudget(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Create an intermediary struct to capture the incoming JSON
	type BudgetInput struct {
		CategoryID uuid.UUID `json:"category_id"`
		Amount     float64   `json:"amount"`
		StartDate  string    `json:"start_date"` // Use string to handle date parsing
		EndDate    string    `json:"end_date"`   // Use string to handle date parsing
	}

	
	var input BudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid input: %v", err), nil, nil)
		return
	}

	// Clean and parse dates before binding to the actual model
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid start_date format: %v", err), nil, nil)
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid end_date format: %v", err), nil, nil)
		return
	}

		// Bind the cleaned data to the Budget model
		newBudget := models.Budget{
			UserID:    userID.(uuid.UUID),
			CategoryID: input.CategoryID,
			Amount:    input.Amount,
			StartDate: startDate,
			EndDate:   endDate,
		}
	
		// Save the new budget
		if err := db.GetDBInstance().Create(&newBudget).Error; err != nil {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to create budget", nil, nil)
			return
		}

	// Send success response
	utils.SendResponse(c, http.StatusCreated, "Budget created successfully", newBudget, nil)
}

// ListBudgets fetches all budgets for a user with optional filters
func ListBudgets(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve query parameters
	period := c.DefaultQuery("period", "current")
	categoryID := c.DefaultQuery("category_id", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	status := c.DefaultQuery("status", "active")

	// Build the query
	var budgets []models.Budget
	query := db.GetDBInstance().Where("user_id = ?", userID)

	// Apply filters
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply period filter based on current, upcoming, or past
	switch period {
	case "current":
		// Example: Retrieve budgets that are ongoing (start_date <= current and end_date >= current)
		query = query.Where("start_date <= ? AND end_date >= ?", time.Now(), time.Now())
	case "upcoming":
		// Example: Retrieve budgets that are yet to start (start_date > current)
		query = query.Where("start_date > ?", time.Now())
	case "past":
		// Example: Retrieve budgets that have already ended (end_date < current)
		query = query.Where("end_date < ?", time.Now())
	default:
		// If the period doesn't match any of the predefined cases, treat it as an invalid request
		utils.SendResponse(c, http.StatusBadRequest, "Invalid period value", nil, nil)
		return
	}

	// Execute the query
	if err := query.Find(&budgets).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch budgets", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Budgets fetched successfully", budgets, nil)
}

// GetSingleBudget fetches detailed information about a specific budget
func GetSingleBudget(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve the budgetId from the URL
	budgetID := c.Param("budgetId")

	// Fetch the budget from the database
	var budget models.Budget
	if err := db.GetDBInstance().Where("user_id = ? AND budget_id = ?", userID, budgetID).First(&budget).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Budget not found", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch budget", nil, nil)
		}
		return
	}

	utils.SendResponse(c, http.StatusOK, "Budget fetched successfully", budget, nil)
}

// UpdateBudget modifies an existing budget
func UpdateBudget(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve the budgetId from the URL
	budgetID := c.Param("budgetId")

	// Fetch the budget from the database
	var budget models.Budget
	if err := db.GetDBInstance().Where("user_id = ? AND budget_id = ?", userID, budgetID).First(&budget).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Budget not found", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch budget", nil, nil)
		}
		return
	}

	// Bind the JSON request data
	var updateData models.Budget
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid input: %v", err), nil, nil)
		return
	}

	// Update fields
	if updateData.Amount != 0 {
		budget.Amount = updateData.Amount
	}
	if !updateData.StartDate.IsZero() {
		budget.StartDate = updateData.StartDate
	}
	if !updateData.EndDate.IsZero() {
		budget.EndDate = updateData.EndDate
	}
	if updateData.CategoryID != uuid.Nil {
		budget.CategoryID = updateData.CategoryID
	}

	// Save the updated budget
	if err := db.GetDBInstance().Save(&budget).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to update budget", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Budget updated successfully", budget, nil)
}

// DeleteBudget removes a specific budget
func DeleteBudget(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve the budgetId from the URL
	budgetID := c.Param("budgetId")

	// Fetch the budget from the database
	var budget models.Budget
	if err := db.GetDBInstance().Where("user_id = ? AND budget_id = ?", userID, budgetID).First(&budget).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Budget not found", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch budget", nil, nil)
		}
		return
	}

	// Delete the budget
	if err := db.GetDBInstance().Delete(&budget).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete budget", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Budget deleted successfully", nil, nil)
}

// BudgetAnalysis compares actual spending against set budgets
func BudgetAnalysis(c *gin.Context) {
	// Get user_id from context
	userID, ok := c.Get("userId")
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found", nil, nil)
		return
	}

	// Retrieve query parameters
	categoryID := c.DefaultQuery("category_id", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// Fetch budgets and expenses within the specified period
	var budgets []models.Budget
	query := db.GetDBInstance().Where("user_id = ?", userID)
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}
	if err := query.Find(&budgets).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch budgets", nil, nil)
		return
	}

	// Calculate total spending for the same period
	var totalSpent float64
	expensesQuery := db.GetDBInstance().Model(&models.Expense{}).Where("user_id = ?", userID)
	if startDate != "" {
		expensesQuery = expensesQuery.Where("date >= ?", startDate)
	}
	if endDate != "" {
		expensesQuery = expensesQuery.Where("date <= ?", endDate)
	}

	// Sum the expenses
	if err := expensesQuery.Select("SUM(amount)").Scan(&totalSpent).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch expenses", nil, nil)
		return
	}

	// Compare budget vs spending
	var analysisResult []struct {
		CategoryID  uuid.UUID `json:"category_id"`
		Category    string    `json:"category"`
		Amount      float64   `json:"amount"`
		TotalSpent  float64   `json:"total_spent"`
		Exceeds     bool      `json:"exceeds"`
	}

	for _, budget := range budgets {
		categoryName := "Unknown" // Fetch category name using the category ID
		if totalSpent > budget.Amount {
			analysisResult = append(analysisResult, struct {
				CategoryID  uuid.UUID `json:"category_id"`
				Category    string    `json:"category"`
				Amount      float64   `json:"amount"`
				TotalSpent  float64   `json:"total_spent"`
				Exceeds     bool      `json:"exceeds"`
			}{
				CategoryID: budget.CategoryID,
				Category:   categoryName,
				Amount:     budget.Amount,
				TotalSpent: totalSpent,
				Exceeds:    true,
			})
		}
	}

	utils.SendResponse(c, http.StatusOK, "Budget analysis fetched successfully", analysisResult, nil)
}
