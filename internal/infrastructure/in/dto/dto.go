package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
}