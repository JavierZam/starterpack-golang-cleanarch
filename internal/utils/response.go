package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// Import os untuk mengecek environment
	"starterpack-golang-cleanarch/internal/utils/errors"
	"starterpack-golang-cleanarch/internal/utils/log"
)

// ErrorResponse struct for consistent error messages returned to clients.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"` // Added omitempty
}

// RespondJSON writes a JSON response to the client with the given status code and data.
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Errorf(context.Background(), "Failed to write JSON response: %v", err)
			http.Error(w, `{"code":"INTERNAL_SERVER_ERROR","message":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// HandleHTTPError maps application errors (AppError interface) to appropriate HTTP responses.
func HandleHTTPError(w http.ResponseWriter, err error, r *http.Request) {
	appErr, ok := err.(errors.AppError) // Try to cast the error to our AppError interface
	if !ok {
		// If it's not an AppError, it's an unexpected internal server error
		log.Errorf(r.Context(), "Unhandled error: %v", err) // Log the actual unhandled error
		RespondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errors.ErrInternalServer.Code(),
			Message: errors.ErrInternalServer.Message(),
			Details: err.Error(), // Include original error for debugging (consider removing/simplifying in production)
		})
		return
	}

	// Log the error based on its HTTP status code for better monitoring
	if appErr.Status() >= http.StatusInternalServerError { // 5xx errors are server errors
		// For server errors (5xx), log the underlying error
		log.Errorf(r.Context(), "Server error occurred: %v", appErr)

		// --- PERUBAHAN AGGRESIF UNTUK DEBUGGING ---
		// Kali ini, KITA AKAN SELALU mengekspos underlying error jika ada,
		// terlepas dari APP_ENV, untuk menemukan akar masalah.
		detailsMessage := appErr.Message() // Default ke pesan AppError
		if unwrappedErr := appErr.Unwrap(); unwrappedErr != nil {
			detailsMessage = unwrappedErr.Error() // Gunakan pesan dari underlying error
		} else {
			detailsMessage = appErr.Error() // Jika tidak ada underlying, gunakan pesan dari AppError.Error()
		}
		// --- AKHIR PERUBAHAN AGGRESIF ---

		RespondJSON(w, appErr.Status(), ErrorResponse{
			Code:    appErr.Code(),
			Message: appErr.Message(),
			Details: detailsMessage, // Sekarang ini akan berisi detail error di development
		})
	} else { // 4xx errors are client errors
		log.Warnf(r.Context(), "Client-side error: %v", appErr)
		RespondJSON(w, appErr.Status(), ErrorResponse{
			Code:    appErr.Code(),
			Message: appErr.Message(),
			Details: formatErrorDetails(appErr.Details()), // Details untuk client errors
		})
	}
}

// formatErrorDetails converts a map of details into a simple string for now.
func formatErrorDetails(details map[string]interface{}) string {
	if len(details) == 0 {
		return ""
	}
	for _, v := range details {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// PaginationRequest is a common struct for handling pagination query parameters.
type PaginationRequest struct {
	Page  int    `json:"page" query:"page" validate:"min=1"`
	Limit int    `json:"limit" query:"limit" validate:"min=1,max=100"`
	Query string `json:"query" query:"query"`
}

// PaginationResponse is a common generic struct for sending paginated data responses.
type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	NextPage   *int  `json:"next_page,omitempty"`
	PrevPage   *int  `json:"prev_page,omitempty"`
}

// ISO8601TimeFormat provides a consistent format for time strings.
const ISO8601TimeFormat = "2006-01-02T15:04:05Z07:00"
