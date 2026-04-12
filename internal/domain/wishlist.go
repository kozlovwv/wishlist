package domain

import "time"

type Wishlist struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	EventDate   time.Time
	PublicToken string
}
