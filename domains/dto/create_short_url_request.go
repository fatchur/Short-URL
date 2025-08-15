package dto

type CreateShortUrlRequest struct {
	LongUrl string `json:"long_url" validate:"required,url"`
}