package dto

type CreateSessionRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	DeviceInfo string `json:"device_info,omitempty"`
	IPAddress  string `json:"ip_address,omitempty"`
}
