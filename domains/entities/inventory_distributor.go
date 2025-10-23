package entities

import (
	"time"
)

type Distributor struct {
	ID          uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Email       string    `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	PhoneNumber *string   `json:"phone_number" gorm:"type:varchar(20)"`
	Address     *string   `json:"address" gorm:"type:text"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`
	CreatedBy   uint      `json:"-" gorm:"index"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy   *uint     `json:"-" gorm:"index"`
}

func (Distributor) TableName() string {
	return "distributors"
}
