package entities

import (
	"time"
)

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	InstitutionID uint      `json:"institution_id" gorm:"not null;index"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	Email         string    `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	PhoneNumber   *string   `json:"phone_number" gorm:"type:varchar(20)"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy     uint      `json:"created_by" gorm:"index"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy     *uint     `json:"updated_by" gorm:"index"`

}

func (User) TableName() string {
	return "users"
}
