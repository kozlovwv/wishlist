package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"wishlist/internal/application"
	"wishlist/internal/infrastructure/in/dto"
	"wishlist/internal/infrastructure/in/response"
)

type AuthHandler struct {
	service  *application.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	token, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.Created(w, dto.AuthResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.OK(w, dto.AuthResponse{Token: token})
}
