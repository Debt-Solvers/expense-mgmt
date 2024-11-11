package db

import (
	"expense-mgmt/internal/models"
	"log"

	"gorm.io/gorm"
)

var defaultExpenseCategories = []models.Category{
	{Name: "Food & Dining", Description: "Restaurants, groceries, and food delivery", ColorCode: "#FFD700"},
	{Name: "Transportation", Description: "Public transport, fuel, car maintenance", ColorCode: "#4682B4"},
	{Name: "Housing", Description: "Rent, mortgage, utilities, maintenance", ColorCode: "#FF6347"},
	{Name: "Entertainment", Description: "Movies, games, hobbies, streaming services", ColorCode: "#8A2BE2"},
	{Name: "Shopping", Description: "Clothing, electronics, personal items", ColorCode: "#32CD32"},
	{Name: "Healthcare", Description: "Medical expenses, medications, insurance", ColorCode: "#FF4500"},
	{Name: "Education", Description: "Tuition, books, courses, training", ColorCode: "#4682B4"},
	{Name: "Utilities", Description: "Electricity, water, internet, phone", ColorCode: "#A9A9A9"},
	{Name: "Travel", Description: "Vacations, business trips, accommodations", ColorCode: "#20B2AA"},
	{Name: "Insurance", Description: "Health, life, car, home insurance", ColorCode: "#FF8C00"},
	{Name: "Others", Description: "Miscellaneous expenses", ColorCode: "#A9A9A9"},
}

// Seed the database by creating default categories..
func SeedDefaultCategories(db *gorm.DB) error {
	// Ensure the categories table exists
	if err := db.AutoMigrate(&models.Category{}); err != nil {
		return err
	}

	// Check if any categories already exist
	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("Categories already seeded, skipping...")
		return nil
	}

	// Seed default categories
	for _, category := range defaultExpenseCategories {
		// Assign a blank user ID for default categories (or adjust as needed)
		category.UserID = nil // or some predefined value if needed
    category.IsDefault = true
		if err := db.Create(&category).Error; err != nil {
			return err
		}
	}

	log.Println("Successfully seeded default categories")
	return nil
}
