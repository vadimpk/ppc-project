package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standardized API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// JSON sends a successful response with data
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: true,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// Error sends an error response
func Error(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorWithCode sends an error response with an error code
func ErrorWithCode(w http.ResponseWriter, statusCode int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Message: message,
			Code:    code,
		},
	}

	json.NewEncoder(w).Encode(response)
}
