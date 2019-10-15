package models

// swagger:model Response
// Response is a generic rest response
type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Paged model
type Paged struct {
	Items  interface{} `json:"items,omitempty"`
	LastId int64       `json:"lastId,omitempty"`
	Limit  int         `json:"limit,omitempty"`
}

// swagger:model LoginResponse
// LoginResponse model
type LoginResponse struct {
	Token     string `json:"token,omitempty"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`
}

// swagger:model LoginRequest
// LoginRequest model
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:model RegisterRequest
// RegisterRequest model
type RegisterRequest struct {
	Name    string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
