package model

type ErrorResponse struct {
	Status  string `json:"status" validate:"required"`
	Message string `json:"message" validate:"required"`
}
