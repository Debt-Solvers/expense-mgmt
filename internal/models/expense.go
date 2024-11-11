package models

import (
	"time"

	"github.com/google/uuid"
)

// Expense represents an individual expense entry associated with a user and category
type Expense struct {
	ExpenseID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"expense_id"`
	UserID             uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID         uuid.UUID     `gorm:"type:uuid;not null" json:"category_id"`
	Amount             float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	Date               time.Time     `gorm:"type:timestamp;not null" json:"date"`
	Description        string        `gorm:"type:text" json:"description"`
	ReceiptID          *uuid.UUID    `gorm:"type:uuid" json:"receipt_id"`
	CreatedAt          time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}