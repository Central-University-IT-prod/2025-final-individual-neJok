package model

type GigaChatTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type GenerateTextRequest struct {
	AdTitle  string  `json:"title" description:"Текст рекламы" binding:"required"`
	Wishes   *string `json:"wishes" binding:"required,max=100"`
	Gender   string  `json:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Location *string `json:"location,omitempty" binding:"omitempty"`
}

type GenerateTextResponse struct {
	Message string `json:"message"`
}

type ModerationRequest struct {
	Status *bool `json:"status" binding:"required"`
}
