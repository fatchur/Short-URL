package entities

import (
	"time"

	"gorm.io/gorm"
)

type UrlSafety struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ShortUrlID uint           `json:"short_url_id" gorm:"not null;uniqueIndex"`
	IsSafe     bool           `json:"is_safe" gorm:"default:false"`
	CheckedAt  time.Time      `json:"checked_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

}