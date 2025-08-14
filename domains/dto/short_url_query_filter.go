package dto

import "time"

type ShortUrlQueryFilter struct {
	UserID    *uint      `json:"user_id,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}
