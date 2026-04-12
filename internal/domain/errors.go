package domain

import "errors"

var (
	ErrWishlistForbidden = errors.New("wishlist access forbidden")

	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
)
