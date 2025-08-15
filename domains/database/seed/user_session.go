package seed

import (
	"time"

	"short-url/domains/entities"
)

var UserSessions = []entities.UserSession{
	{
		ID:           1,
		UserID:       1,
		SessionToken: "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
		SecretKey:    "secret1234567890abcd1234567890abcd1234567890abcd1234567890abcd12",
		DeviceInfo:   StringPtr("Chrome Browser / MacOS"),
		IPAddress:    StringPtr("127.0.0.1"),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	{
		ID:           2,
		UserID:       2,
		SessionToken: "efgh5678901234efgh5678901234efgh5678901234efgh5678901234efgh5678",
		SecretKey:    "secret5678901234efgh5678901234efgh5678901234efgh5678901234efgh56",
		DeviceInfo:   StringPtr("Firefox Browser / Windows"),
		IPAddress:    StringPtr("127.0.0.1"),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	{
		ID:           3,
		UserID:       3,
		SessionToken: "ijkl9012345678ijkl9012345678ijkl9012345678ijkl9012345678ijkl9012",
		SecretKey:    "secret9012345678ijkl9012345678ijkl9012345678ijkl9012345678ijkl90",
		DeviceInfo:   StringPtr("Safari Browser / iOS"),
		IPAddress:    StringPtr("127.0.0.1"),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
}