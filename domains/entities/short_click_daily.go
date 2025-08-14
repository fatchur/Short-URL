package entities

import (
	"time"

	"gorm.io/gorm"
)

type ShortClickDaily struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ShortUrlID uint           `json:"short_url_id" gorm:"not null"`
	Date       time.Time      `json:"date" gorm:"type:date;not null"`
	NumRequest int            `json:"num_request" gorm:"default:0"`
	CreatedAt  time.Time      `json:"created_at"`
	CreatedBy  uint           `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at"`
	UpdatedBy  uint           `json:"updated_by"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships can be loaded separately to avoid circular references
}
