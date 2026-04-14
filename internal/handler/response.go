package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    any    `json:"message"`
	Error      string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func SendError(w http.ResponseWriter, status int, message any, errName string) {
	JSON(w, status, ErrorResponse{
		StatusCode: status,
		Message:    message,
		Error:      errName,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	SendError(w, http.StatusBadRequest, message, "Bad Request")
}

func Unauthorized(w http.ResponseWriter, message string) {
	SendError(w, http.StatusUnauthorized, message, "Unauthorized")
}

func Forbidden(w http.ResponseWriter, message string) {
	SendError(w, http.StatusForbidden, message, "Forbidden")
}

func NotFound(w http.ResponseWriter, message string) {
	SendError(w, http.StatusNotFound, message, "Not Found")
}

func Conflict(w http.ResponseWriter, message string) {
	SendError(w, http.StatusConflict, message, "Conflict")
}

func InternalError(w http.ResponseWriter, message string) {
	SendError(w, http.StatusInternalServerError, message, "Internal Server Error")
}
