package models

import (
	"time"

	"gorm.io/gorm"
)

// Budget represents a spending cap set by a user on a specific category
type Budget struct {
	ID         string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"budget_id"`
	UserID     string         `gorm:"type:uuid;not null" json:"user_id"`           // Foreign key to the User
	CategoryID *string        `gorm:"type:uuid" json:"category_id"`                // Foreign key to the Category (nullable)
	Amount     float64        `gorm:"type:decimal(10,2);check:amount >= 0;not null" json:"amount"` // Budget amount
	StartDate  time.Time      `gorm:"type:date;not null" json:"start_date"`        // Budget start date
	EndDate    time.Time      `gorm:"type:date;not null" json:"end_date"`          // Budget end date
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}