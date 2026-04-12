package dto

import "time"

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

type WishlistRequest struct {
	Title       string `json:"title"       validate:"required,min=1,max=255"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"  validate:"required"`
}

type WishlistResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	PublicToken string    `json:"public_token"`
}

type ItemRequest struct {
	Title       string `json:"title"       validate:"required,min=1,max=255"`
	Description string `json:"description"`
	URL         string `json:"url"         validate:"omitempty,url"`
	Priority    int    `json:"priority"    validate:"required,min=1,max=10"`
}

type ItemResponse struct {
	ID          int64  `json:"id"`
	WishlistID  int64  `json:"wishlist_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Priority    int    `json:"priority"`
	IsReserved  bool   `json:"is_reserved"`
}
