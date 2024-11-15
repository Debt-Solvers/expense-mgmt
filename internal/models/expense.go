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



// ExpenseAnalysisResult represents the result of the expense analysis query.
type ExpenseAnalysisResult struct {
	Period            string             `json:"period"`
	TotalSpending     float64            `json:"total_spending"`
	AverageSpending   float64            `json:"average_spending"`
	CategoryBreakdown map[string]float64 `json:"category_breakdown,omitempty"`
	HighestExpense    float64            `json:"highest_expense,omitempty"`
	RecurringTotal    float64            `json:"recurring_total,omitempty"`
}
