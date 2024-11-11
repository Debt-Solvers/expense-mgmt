package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents a user-defined or default spending category
type Category struct {
	ID 					uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"category_id"`
	UserID      *uuid.UUID     `gorm:"type:uuid" json:"user_id,omitempty"`           // Nullable for default categories
	Name        string         `gorm:"size:50;not null" json:"name"`                 // e.g., "Food", "Utilities"
	Description string         `gorm:"type:text" json:"description"`                 // Optional description
	ColorCode   string         `gorm:"size:7" json:"color_code"`                     // Optional color code (e.g., #FFFFFF)
	IsDefault   bool           `gorm:"default:false" json:"is_default"`              // True if the category is default
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`            // Soft delete
}

// FetchDefaultCategories retrieves all categories marked as default from the database.
func FetchDefaultCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	err := db.Where("is_default = ?", true).Find(&categories).Error
	return categories, err
}

