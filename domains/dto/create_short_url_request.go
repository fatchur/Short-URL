package dto

type CreateShortUrlRequest struct {
	LongUrl string `json:"long_url" validate:"required,url"`
	UserID  uint   `json:"user_id" validate:"required"`
}