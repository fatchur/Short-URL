package entities

import (
	"time"

	"gorm.io/gorm"
)

type ShortUrl struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	LongUrl   string         `json:"long_url" gorm:"type:text;not null"`
	ShortCode string         `json:"short_code" gorm:"type:varchar(10);uniqueIndex;not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	ExpireAt  *time.Time     `json:"expire_at"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy uint           `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy uint           `json:"updated_by"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	User             User              `json:"user" gorm:"foreignKey:UserID"`
	ShortClickDailys []ShortClickDaily `json:"short_click_dailys" gorm:"foreignKey:ShortUrlID"`
}
