package controller

import (
	"expense-mgmt/internal/models"
	"expense-mgmt/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDefaultCategories returns a list of pre-defined default categories.
func GetDefaultCategories(c *gin.Context) {
	// Call the model function to fetch default categories
	categories, err := models.FetchDefaultCategories() // Directly fetch from model, no need to bind JSON
	if err != nil {
		// If an error occurs while fetching, send an error response
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch default categories", nil, gin.H{"error": err.Error()})
		return
	}

	// Send a success response with the list of default categories
	utils.SendResponse(c, http.StatusOK, "Fetched default categories successfully", gin.H{"categories": categories}, nil)
}


// ListUserCategories lists all categories for the authenticated user.
func ListUserCategories(c *gin.Context) {
	// Fetch query params for sorting, filtering, pagination, etc.
	categories, err := models.FetchUserCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory allows users to add new custom categories.
func CreateCategory(c *gin.Context) {
	var categoryRequest CategoryRequest
	if err := c.ShouldBindJSON(&categoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	category, err := models.CreateUserCategory(c, categoryRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"category": category})
}

// GetCategoryDetails returns information for a specific category.
func GetCategoryDetails(c *gin.Context) {
	categoryId := c.Param("categoryId")
	category, err := models.FetchCategoryDetails(categoryId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category})
}

// UpdateCategory updates an existing category's details.
func UpdateCategory(c *gin.Context) {
	var categoryUpdate CategoryUpdateRequest
	categoryId := c.Param("categoryId")
	if err := c.ShouldBindJSON(&categoryUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}
	err := models.UpdateUserCategory(categoryId, categoryUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}
	c.Status(http.StatusOK)
}

// GetCategorySummary returns usage summary for a specific category.
func GetCategorySummary(c *gin.Context) {
	categoryId := c.Param("categoryId")
	summary, err := models.GetCategoryUsageSummary(categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category summary"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

// GetCategoryBudgetStatus returns budget status for a specific category.
func GetCategoryBudgetStatus(c *gin.Context) {
	categoryId := c.Param("categoryId")
	budgetStatus, err := models.GetCategoryBudgetStatus(categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budget status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"budgetStatus": budgetStatus})
}