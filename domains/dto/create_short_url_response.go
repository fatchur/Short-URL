package dto

type CreateShortUrlResponse struct {
	ID        uint   `json:"id"`
	ShortCode string `json:"short_code"`
	LongUrl   string `json:"long_url"`
	UserID    uint   `json:"user_id"`
}