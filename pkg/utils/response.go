package utils

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": message})
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RespondWithValidationError sends validation error response
func RespondWithValidationError(w http.ResponseWriter, statusCode int, errors []ValidationError) {
	RespondWithJSON(w, statusCode, map[string]interface{}{
		"error":   "Validation failed",
		"details": errors,
	})
}

func GetUserClaimsFromContext(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value("userClaims").(jwt.MapClaims)
	return claims, ok
}
