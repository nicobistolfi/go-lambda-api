package models

type ErrorResponse struct {
	Error string `json:"error" example:"Missing API key"`
}
