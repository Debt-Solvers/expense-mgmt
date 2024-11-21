package controller

import (
	"errors"
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

	// Intermediary struct to capture incoming JSON
	type BudgetInput struct {
		CategoryID uuid.UUID `json:"category_id" binding:"required"` // Ensure category_id is provided
		Amount     float64   `json:"amount" binding:"required,gt=0"` // Ensure amount is positive
		StartDate  string    `json:"start_date" binding:"required"`  // Ensure start_date is provided
		EndDate    string    `json:"end_date" binding:"required"`    // Ensure end_date is provided
	}

	var input BudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid input: %v", err), nil, nil)
		return
	}

	// Parse dates and validate their format
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid start_date format (expected YYYY-MM-DD): %v", err), nil, nil)
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid end_date format (expected YYYY-MM-DD): %v", err), nil, nil)
		return
	}

	// Ensure end_date is not earlier than start_date
	if endDate.Before(startDate) {
		utils.SendResponse(c, http.StatusBadRequest, "end_date must be later than or equal to start_date", nil, nil)
		return
	}

	// Ensure the budget does not overlap with existing budgets for the same category
	var overlappingBudgetExists bool
	err = db.GetDBInstance().Model(&models.Budget{}).
		Where("user_id = ? AND category_id = ? AND NOT (end_date < ? OR start_date > ?)", userID, input.CategoryID, startDate, endDate).
		Select("COUNT(1) > 0").
		Scan(&overlappingBudgetExists).Error
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to check for overlapping budgets", nil, nil)
		return
	}

	if overlappingBudgetExists {
		utils.SendResponse(c, http.StatusBadRequest, "Budget period overlaps with an existing budget for the same category", nil, nil)
		return
	}

	// Create the new budget
	newBudget := models.Budget{
		UserID:     userID.(uuid.UUID),
		CategoryID: input.CategoryID,
		Amount:     input.Amount,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	// Save the budget to the database
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
	categoryID := c.DefaultQuery("category_id", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

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
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found in context", nil, nil)
		return
	}

	// Retrieve the budgetId from the URL
	budgetID := c.Param("budgetId")
	if budgetID == "" {
		utils.SendResponse(c, http.StatusBadRequest, "Budget ID is required", nil, nil)
		return
	}

	// Validate the budgetID format
	if _, err := uuid.Parse(budgetID); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid budget ID format", nil, nil)
		return
	}

	// Fetch the budget from the database
	var budget models.Budget
	err := db.GetDBInstance().Where("user_id = ? AND budget_id = ?", userID, budgetID).First(&budget).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendResponse(c, http.StatusNotFound, "Budget not found", nil, nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, "Error retrieving budget from database", nil, nil)
		}
		return
	}

	// Bind the JSON request data
	var updateData struct {
		Amount     *float64   `json:"amount"`      // Use pointers to distinguish between unset and zero values
		StartDate  *string    `json:"start_date"`  // Use string to validate and parse date later
		EndDate    *string    `json:"end_date"`    // Use string for date validation
		CategoryID *uuid.UUID `json:"category_id"` // Optional category update
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid input: %v", err), nil, nil)
		return
	}

	// Validate and update each field
	if updateData.Amount != nil {
		if *updateData.Amount <= 0 {
			utils.SendResponse(c, http.StatusBadRequest, "Amount must be greater than zero", nil, nil)
			return
		}
		budget.Amount = *updateData.Amount
	}

	if updateData.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *updateData.StartDate)
		if err != nil {
			utils.SendResponse(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD", nil, nil)
			return
		}
		budget.StartDate = startDate
	}

	if updateData.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *updateData.EndDate)
		if err != nil {
			utils.SendResponse(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD", nil, nil)
			return
		}
		if updateData.StartDate != nil {
			startDate, _ := time.Parse("2006-01-02", *updateData.StartDate)
			if endDate.Before(startDate) {
				utils.SendResponse(c, http.StatusBadRequest, "End date cannot be earlier than start date", nil, nil)
				return
			}
		} else if endDate.Before(budget.StartDate) {
			utils.SendResponse(c, http.StatusBadRequest, "End date cannot be earlier than the current start date", nil, nil)
			return
		}
		budget.EndDate = endDate
	}

	if updateData.CategoryID != nil {
		budget.CategoryID = *updateData.CategoryID
	}

	// Save the updated budget
	if err := db.GetDBInstance().Save(&budget).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to update budget", nil, nil)
		return
	}

	// Respond with the updated budget
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
		utils.SendResponse(c, http.StatusUnauthorized, "User ID not found in context", nil, nil)
		return
	}

	// Retrieve query parameters
	categoryID := c.DefaultQuery("category_id", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// Validate date formats if provided
	var err error
	if startDate != "" {
		if _, err = time.Parse("2006-01-02", startDate); err != nil {
			utils.SendResponse(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD", nil, nil)
			return
		}
	}
	if endDate != "" {
		if _, err = time.Parse("2006-01-02", endDate); err != nil {
			utils.SendResponse(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD", nil, nil)
			return
		}
	}

	// Fetch budgets for the specified period
	var budgets []models.Budget
	budgetQuery := db.GetDBInstance().Where("user_id = ?", userID)

	if categoryID != "" {
		budgetQuery = budgetQuery.Where("category_id = ?", categoryID)
	}
	if startDate != "" {
		budgetQuery = budgetQuery.Where("start_date >= ?", startDate)
	}
	if endDate != "" {
		budgetQuery = budgetQuery.Where("end_date <= ?", endDate)
	}
	if err := budgetQuery.Find(&budgets).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch budgets", nil, nil)
		return
	}

	if len(budgets) == 0 {
		utils.SendResponse(c, http.StatusNotFound, "No budgets found for the specified period", nil, nil)
		return
	}

	// Map to hold total spending per category
	categorySpendMap := make(map[uuid.UUID]float64)

	// Fetch total spending per category for the specified period
	expensesQuery := db.GetDBInstance().Model(&models.Expense{}).Where("user_id = ?", userID)
	if startDate != "" {
		expensesQuery = expensesQuery.Where("date >= ?", startDate)
	}
	if endDate != "" {
		expensesQuery = expensesQuery.Where("date <= ?", endDate)
	}
	if categoryID != "" {
		expensesQuery = expensesQuery.Where("category_id = ?", categoryID)
	}

	rows, err := expensesQuery.Select("category_id, SUM(amount) as total_spent").Group("category_id").Rows()
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch expense data", nil, nil)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var catID uuid.UUID
		var totalSpent float64
		if err := rows.Scan(&catID, &totalSpent); err != nil {
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to process expense data", nil, nil)
			return
		}
		categorySpendMap[catID] = totalSpent
	}

	// Prepare analysis results
	type AnalysisResult struct {
		CategoryID   uuid.UUID `json:"category_id"`
		Category     string    `json:"category"`
		Amount       float64   `json:"budgeted_amount"`
		TotalSpent   float64   `json:"total_spent"`
		Remaining    float64   `json:"remaining_budget"`
		Percentage   float64   `json:"percentage_spent"`
		Exceeds      bool      `json:"exceeds_budget"`
	}

	var analysisResults []AnalysisResult

	for _, budget := range budgets {
		totalSpent := categorySpendMap[budget.CategoryID]
		remaining := budget.Amount - totalSpent
		percentage := 0.0
		if budget.Amount > 0 {
			percentage = (totalSpent / budget.Amount) * 100
		}
		exceeds := totalSpent > budget.Amount

		// Fetch category name (use "Unknown" if unavailable)
		categoryName := "Unknown" // Default to Unknown in case of errors

		var category struct {
      Name string
    }
 
    err := db.GetDBInstance().
        Table("categories").
        Select("name").
        Where("id = ?", budget.CategoryID).
        Scan(&category).Error
    if err != nil {
      // Log the error for debugging, but don't fail the entire process
      fmt.Printf("Failed to fetch category name for category_id: %s, Error: %v\n", budget.CategoryID, err)
    }else {
      categoryName = category.Name // Assign the fetched name
    }
	
		analysisResults = append(analysisResults, AnalysisResult{
			CategoryID:   budget.CategoryID,
			Category:     categoryName,
			Amount:       budget.Amount,
			TotalSpent:   totalSpent,
			Remaining:    remaining,
			Percentage:   percentage,
			Exceeds:      exceeds,
		})
	}

	// Send the analysis results
	utils.SendResponse(c, http.StatusOK, "Budget analysis fetched successfully", analysisResults, nil)
}
