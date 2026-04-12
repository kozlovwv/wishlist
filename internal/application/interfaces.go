package application

import (
	"context"
	"wishlist/internal/domain"
)

type WishlistRepository interface {
	Create(ctx context.Context, wishlist domain.Wishlist) (domain.Wishlist, error)
	FindByID(ctx context.Context, id int64) (domain.Wishlist, error)
	FindByUserID(ctx context.Context, userID int64) ([]domain.Wishlist, error)
	FindByPublicToken(ctx context.Context, token string) (domain.Wishlist, error)
	Update(ctx context.Context, wishlist domain.Wishlist) (domain.Wishlist, error)
	Delete(ctx context.Context, id int64) error
}

type ItemRepository interface {
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	FindByID(ctx context.Context, id int64) (domain.Item, error)
	FindByWishlistID(ctx context.Context, wishlistID int64) ([]domain.Item, error)
	Update(ctx context.Context, item domain.Item) (domain.Item, error)
	Delete(ctx context.Context, id int64) error
	Reserve(ctx context.Context, id int64) error
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type PublicTokenGenerator interface {
	Generate() (string, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
 
type TokenManager interface {
	Generate(userID int64) (string, error)
	Parse(token string) (int64, error)
}