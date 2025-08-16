package seed

import (
	"time"

	"short-url/domains/entities"
)

var UserSessions = []entities.UserSession{
	{
		UserID:      1,
		SessionCode: "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
		SecretKey:   "secret1234567890abcd1234567890abcd1234567890abcd1234567890abcd12",
		DeviceInfo:  StringPtr("Chrome Browser / MacOS"),
		IPAddress:   StringPtr("127.0.0.1"),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		UserID:      2,
		SessionCode: "efgh5678901234efgh5678901234efgh5678901234efgh5678901234efgh5678",
		SecretKey:   "secret5678901234efgh5678901234efgh5678901234efgh5678901234efgh56",
		DeviceInfo:  StringPtr("Firefox Browser / Windows"),
		IPAddress:   StringPtr("127.0.0.1"),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		UserID:      3,
		SessionCode: "ijkl9012345678ijkl9012345678ijkl9012345678ijkl9012345678ijkl9012",
		SecretKey:   "secret9012345678ijkl9012345678ijkl9012345678ijkl9012345678ijkl90",
		DeviceInfo:  StringPtr("Safari Browser / iOS"),
		IPAddress:   StringPtr("127.0.0.1"),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}