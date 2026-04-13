package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"wishlist/internal/domain"
	"wishlist/internal/infrastructure/in/dto"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: msg})
}

func Unauthorized(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: msg})
}

func Forbidden(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusForbidden, dto.ErrorResponse{Error: msg})
}

func NotFound(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusNotFound, dto.ErrorResponse{Error: msg})
}

func Conflict(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusConflict, dto.ErrorResponse{Error: msg})
}

func InternalError(w http.ResponseWriter) {
	JSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
}

func FromDomainError(w http.ResponseWriter, err error) {
	switch {
	case isAny(err, domain.ErrUserNotFound, domain.ErrWishlistNotFound, domain.ErrItemNotFound):
		NotFound(w, err.Error())
	case isAny(err, domain.ErrUserAlreadyExists, domain.ErrItemReserved):
		Conflict(w, err.Error())
	case isAny(err, domain.ErrInvalidPassword, domain.ErrInvalidToken):
		Unauthorized(w, err.Error())
	case isAny(err, domain.ErrWishlistForbidden):
		Forbidden(w, err.Error())
	default:
		InternalError(w)
	}
}

func isAny(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
