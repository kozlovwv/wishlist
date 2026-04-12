package domain

import "errors"

var (
	ErrWishlistForbidden = errors.New("wishlist access forbidden")
	ErrWishlistNotFound  = errors.New("wishlist not found")
	ErrItemForbidden     = errors.New("item access forbidden")
	ErrItemNotFound      = errors.New("item not found")
	ErrItemReserved      = errors.New("item is already reserved")

	ErrInvalidToken = errors.New("invalid token")

	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
)
