package entities

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	SessionToken string         `json:"session_token" gorm:"type:varchar(255);uniqueIndex;not null"`
	SecretKey    string         `json:"secret_key" gorm:"type:varchar(255);not null"`
	DeviceInfo   *string        `json:"device_info" gorm:"type:text"`
	IPAddress    *string        `json:"ip_address" gorm:"type:varchar(45)"`
	ExpiresAt    time.Time      `json:"expires_at" gorm:"not null"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}