package domain

type RegisterDTO struct {
	Name     string   `json:"name" example:"Иван Иванов"`
	Email    string   `json:"email" example:"ivan@example.com"`
	Password string   `json:"password" example:"secret123"`
	Roles    []string `json:"roles,omitempty"`
}

type LoginDTO struct {
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"secret123"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type JWTPayload struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode" example:"400"`
	Message    string `json:"message" example:"Bad request message"`
	Error      string `json:"error" example:"Bad Request"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
