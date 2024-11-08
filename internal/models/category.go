package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a user-defined or default spending category
type Category struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"category_id"`
	UserID      string         `gorm:"type:uuid;not null" json:"user_id"`            // Foreign key to the User
	Name        string         `gorm:"size:50;not null" json:"name"`                 // e.g., "Food", "Utilities"
	Description string         `gorm:"type:text" json:"description"`                 // Optional description
	ColorCode   string         `gorm:"size:7" json:"color_code"`                    // Optional color code (e.g., #FFFFFF)
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}

