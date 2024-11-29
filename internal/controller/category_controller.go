package controller

import (
	"errors"
	"expense-mgmt/db"
	"expense-mgmt/internal/common"
	"expense-mgmt/internal/models"
	"expense-mgmt/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetDefaultCategories handles fetching the list of default categories
func GetDefaultCategories(context *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Call the model function to fetch default categories
	categories, err := models.FetchDefaultCategories(DB)
	if err != nil {
		// If fetching fails, send a standardized error response
		utils.SendResponse(context, http.StatusInternalServerError, "Failed to fetch default categories", nil, gin.H{"error": err.Error()})
		return
	}

	// Send a success response with the list of default categories
	utils.SendResponse(context, http.StatusOK, "Fetched default categories successfully", gin.H{"categories": categories}, nil)
}

// CreateCustomCategory handles the creation of a custom category for a user
func CreateCustomCategory(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Bind JSON input to Category struct
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category data"})
		return
	}

	// Get the user ID from the JWT token in the middleware (assumed to be set)
	userID, exists := c.Get("userId")
	if !exists {
		utils.SendResponse(c, http.StatusUnauthorized, "User not authorized", nil, nil)
		return
	}
	
	// Initialize variable for parsed UUID
	var parsedUserID uuid.UUID

	// Check the type of userID and convert if necessary
	switch uid := userID.(type) {
	case uuid.UUID:
			parsedUserID = uid // If it's already a UUID, use it directly
	case string:
			var err error
			parsedUserID, err = uuid.Parse(uid) // Parse if it's a string
			if err != nil {
					utils.SendResponse(c, http.StatusBadRequest, "Invalid user ID format", nil, nil)
					return
			}
	default:
			utils.SendResponse(c, http.StatusUnauthorized, "User not authorized", nil, nil)
			return
	}

	// Set IsDefault to false for user-defined categories
	category.IsDefault = false
	category.UserID = &parsedUserID

	// Basic validation (check if essential fields are provided)
	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
		return
	}

	// Trim the category name to remove any leading/trailing spaces
	category.Name = strings.TrimSpace(category.Name)

	 // Check if the category name already exists (case-insensitive check)
	 var existingCategory models.Category
	 if err := DB.Where("user_id = ? AND LOWER(name) = LOWER(?)", parsedUserID, category.Name).First(&existingCategory).Error; err == nil {
		// If a record is found, it means the category name already exists
			c.JSON(http.StatusConflict, gin.H{"error": "Category name already exists"})
			return
	 }

	// If ColorCode is empty, generate a random color
	if category.ColorCode == "" {
		category.ColorCode = common.GenerateRandomColor()
	}
	
	// Insert the category directly into the database
	if err := DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create custom category"})
		return
	}

	// Return a success response with the created category ID
	utils.SendResponse(c, http.StatusCreated, "Category created successfully", category.ID, nil)
}

// ListUserCategories handles fetching all categories for the authenticated user
func ListUserCategories(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Get the user ID from the middleware context (assuming user ID is added in middleware)
	userID, exists := c.Get("userId")
	if !exists {
		utils.SendResponse(c, http.StatusUnauthorized, "User not authorized", nil, nil)
		return
	}

	// Ensure userID is of the correct type (UUID)
	_, ok := userID.(uuid.UUID)
	if !ok {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid user ID format", nil, nil)
		return
	}

	// Fetch default categories
	var defaultCategories []models.Category
	if err := DB.Where("is_default = ?", true).
		Order("name ASC"). // sort default categories in ASC order
		Find(&defaultCategories).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Error retrieving default categories", nil, nil)
		return
	}

	// Fetch custom categories for the user
	var customCategories []models.Category
	if err := DB.Where("user_id = ?", userID).
		Order("name ASC"). // Sort in ASC order
		Find(&customCategories).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Error retrieving custom categories", nil, nil)
		return
	}

		// Merge default and custom categories
    allCategories := append(defaultCategories, customCategories...)

	// Return the result as a response
	utils.SendResponse(c, http.StatusOK, "Categories retrieved successfully", allCategories, nil)
}


// GetCategoryDetails handles fetching the details of a specific category by its ID
func GetCategoryDetails(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Get the category ID from the request parameters
	categoryID := c.Param("categoryId")

	// Check if categoryId is not empty
	if categoryID == "" {
		utils.SendResponse(c, http.StatusBadRequest, "Category ID is required", nil, nil)
		return
	}

	// Try parsing the categoryId into UUID
	parsedCategoryId, err := uuid.Parse(categoryID)
	if err != nil {
			utils.SendResponse(c, http.StatusBadRequest, "Invalid category ID format", nil, nil)
			return
	}

	// Fetch the category from the database
	var category models.Category
	if err := DB.First(&category, "id = ?", parsedCategoryId).Error; err != nil {
			utils.SendResponse(c, http.StatusNotFound, "Category not found", nil, nil)
			return
	}

	// Return the category details in the response
	utils.SendResponse(c, http.StatusOK, "Category details fetched successfully", category, nil)
}


// UpdateCustomCategory allows the user to update a custom category they created
func UpdateCategory(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Retrieve the category ID from the URL parameter
	categoryID := c.Param("categoryId")
	fmt.Println("Received category ID:", categoryID)

	// Retrieve the user ID from the token or context
	userID := c.MustGet("userId") //.(uuid.UUID)

	// Check if categoryId is not empty
	if categoryID == "" {
		utils.SendResponse(c, http.StatusBadRequest, "Category ID is required", nil, nil)
		return
	}

	// Try parsing the categoryId into UUID
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid category ID format", nil, nil)
		return
	}

	// Retrieve the category by ID and user ID, ensuring it’s a custom category
	var category models.Category
	if err := DB.Where("id = ? AND user_id = ? AND is_default = ?", parsedCategoryID, userID, false).
			First(&category).Error; err != nil {
			// Differentiate errors: category not found or user cannot update a default category
			if err == gorm.ErrRecordNotFound {
				utils.SendResponse(c, http.StatusForbidden, "Cannot update this category", nil, nil)
				return
			}
			utils.SendResponse(c, http.StatusInternalServerError, "Error fetching category", nil, nil)
			return
	}

	// Bind the updated fields from the request body
	var updatedCategoryData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ColorCode   string `json:"color_code"`
	}
	if err := c.ShouldBindJSON(&updatedCategoryData); err != nil {
			utils.SendResponse(c, http.StatusBadRequest, "Invalid request body", nil, nil)
			return
	}

	// Update the fields if provided
	if updatedCategoryData.Name != "" {
		category.Name = updatedCategoryData.Name
	}
	if updatedCategoryData.Description != "" {
		category.Description = updatedCategoryData.Description
	}

	// Update ColorCode: use provided value, or keep existing if available, or generate new if not
	if updatedCategoryData.ColorCode != "" {
			category.ColorCode = updatedCategoryData.ColorCode
	} else if category.ColorCode == "" {
			category.ColorCode = common.GenerateRandomColor()
	}

	// Save the updated category to the database
	if err := DB.Save(&category).Error; err != nil {
			utils.SendResponse(c, http.StatusInternalServerError, "Error updating category", nil, nil)
			return
	}

	// Send the response with the updated category data
	utils.SendResponse(c, http.StatusOK, "Category updated successfully", category, nil)
}

// DeleteCustomCategory handles the deletion of a specific custom category by its ID
func DeleteCategory(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Get the category ID from the request parameters
	categoryID := c.Param("categoryId")

	// Check if categoryID is provided
	if categoryID == "" {
		utils.SendResponse(c, http.StatusBadRequest, "Category ID is required", nil, nil)
		return
	}

	// Parse the category ID into a UUID
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid category ID format", nil, nil)
		return
	}

	// Get the user ID from the context (assuming middleware has set this for authenticated users)
	userID := c.MustGet("userId")

	// Find the category by ID and ensure it belongs to the user and is not a default category
	var category models.Category
	if err := DB.Where("id = ? AND user_id = ? AND is_default = ?", parsedCategoryID, userID, false).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendResponse(c, http.StatusNotFound, "Category not found or not allowed to delete", nil, nil)
				return
			}
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to fetch category", nil, nil)
			return
	}

	// Permanently delete the category from the database
	if err := DB.Unscoped().Delete(&category).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete category", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Category deleted successfully", nil, nil)
}

