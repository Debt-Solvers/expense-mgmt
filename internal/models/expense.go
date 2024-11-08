package models

import (
	"time"

	// "github.com/Debt-Solvers/BE-auth-service/blob/main/internal/models/user.go" // this line will be used to import user model from github
	"gorm.io/gorm"
)

// Expense represents an individual expense entry associated with a user and category
type Expense struct {
	ID                  string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"expense_id"`
	UserID              string         `gorm:"type:uuid;not null" json:"user_id"`             // Foreign key to the User
	CategoryID          *string        `gorm:"type:uuid" json:"category_id"`                  // Foreign key to the Category (nullable)
	Amount              float64        `gorm:"type:decimal(10,2);not null" json:"amount"`     // Amount spent
	Date                time.Time      `gorm:"type:timestamp with time zone;not null" json:"date"`
	Description         string         `gorm:"type:text" json:"description"`                  // Optional description or notes
	ReceiptID           *string        `gorm:"type:uuid" json:"receipt_id"`                   // Foreign key to the Receipt (nullable)
	IsRecurring         bool           `gorm:"default:false" json:"is_recurring"`
	RecurringIntervalDays int          `gorm:"type:int;check:recurring_interval_days >= 0" json:"recurring_interval_days,omitempty"` // Optional interval in days
	UpdatedAt  					time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  					gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}