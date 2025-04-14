package model

type TimeSetRequest struct {
	CurrentDate *int32 `json:"current_date" binding:"required,min=0"`
}
