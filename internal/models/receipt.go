package models

import (
	"time"

	"gorm.io/gorm"
)

// Receipt represents a receipt linked to a specific expense, with OCR and image metadata
type Receipt struct {
	ID         string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"receipt_id"`
	ImageURL   string         `gorm:"type:text;not null" json:"image_url"`              // URL or path to the stored image
	OCRData    string         `gorm:"type:text" json:"ocr_data"`                        // Text extracted from the receipt via OCR
	ScannedDate time.Time     `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"scanned_date"`
	ExpenseID  *string        `gorm:"type:uuid" json:"expense_id,omitempty"`            // Foreign key to the Expense (nullable)
	Date  			time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}