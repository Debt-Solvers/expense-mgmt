package db

import (
	"expense-mgmt/internal/models"
	"log"

	"gorm.io/gorm"
)

// var defaultExpenseCategories = []models.Category{
// 	{Name: "Food & Dining", Description: "Restaurants, groceries, and food delivery", Type: "expense"},
// 	{Name: "Transportation", Description: "Public transport, fuel, car maintenance", Type: "expense"},
// 	{Name: "Housing", Description: "Rent, mortgage, utilities, maintenance", Type: "expense"},
// 	{Name: "Entertainment", Description: "Movies, games, hobbies, streaming services", Type: "expense"},
// 	{Name: "Shopping", Description: "Clothing, electronics, personal items", Type: "expense"},
// 	{Name: "Healthcare", Description: "Medical expenses, medications, insurance", Type: "expense"},
// 	{Name: "Education", Description: "Tuition, books, courses, training", Type: "expense"},
// 	{Name: "Utilities", Description: "Electricity, water, internet, phone", Type: "expense"},
// 	{Name: "Travel", Description: "Vacations, business trips, accommodations", Type: "expense"},
// 	{Name: "Insurance", Description: "Health, life, car, home insurance", Type: "expense"},
// }

// var defaultIncomeCategories = []models.Category{
// 	{Name: "Salary", Description: "Regular employment income", Type: "income"},
// 	{Name: "Freelance", Description: "Independent contractor earnings", Type: "income"},
// 	{Name: "Investments", Description: "Dividends, interest, capital gains", Type: "income"},
// 	{Name: "Gifts", Description: "Money received as gifts", Type: "income"},
// 	{Name: "Rental", Description: "Income from property rentals", Type: "income"},
// }

func SeedDefaultCategories(db *gorm.DB) error {
	// Create categories table if it doesn't exist
	err := db.AutoMigrate(&models.Category{})
	if err != nil {
			return err
	}

	// Check if categories already exist
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
			log.Println("Categories already seeded, skipping...")
			return nil
	}

	// Create default expense categories
	for _, category := range defaultExpenseCategories {
			if err := db.Create(&category).Error; err != nil {
					return err
			}
	}

	// Create default income categories
	// for _, category := range defaultIncomeCategories {
	// 		if err := db.Create(&category).Error; err != nil {
	// 				return err
	// 		}
	// }

	log.Println("Successfully seeded default categories")
	return nil
}
