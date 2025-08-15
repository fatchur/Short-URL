package dto

type BaseResponse struct {
	Success    bool        `json:"success"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	APIVersion string      `json:"api_version"`
	Data       interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(status int, message string, data interface{}) BaseResponse {
	return BaseResponse{
		Success:    true,
		Status:     status,
		Message:    message,
		APIVersion: "v1",
		Data:       data,
	}
}

func NewErrorResponse(status int, message string) BaseResponse {
	return BaseResponse{
		Success:    false,
		Status:     status,
		Message:    message,
		APIVersion: "v1",
		Data:       nil,
	}
}