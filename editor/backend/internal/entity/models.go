package entity

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
	}
}
