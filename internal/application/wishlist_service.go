package application

import (
	"context"
	"fmt"
	"wishlist/internal/domain"
)

type WishlistService struct {
	wishlistRepo   WishlistRepository
	tokenGenerator PublicTokenGenerator
}

func NewWishlistService(wishlistRepo WishlistRepository, tokenGenerator PublicTokenGenerator) *WishlistService {
	return &WishlistService{
		wishlistRepo:   wishlistRepo,
		tokenGenerator: tokenGenerator,
	}
}

func (s *WishlistService) Create(ctx context.Context, w domain.Wishlist) (domain.Wishlist, error) {
	token, err := s.tokenGenerator.Generate()
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("generate public token: %w", err)
	}

	w.PublicToken = token

	created, err := s.wishlistRepo.Create(ctx, w)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("create wishlist: %w", err)
	}
	return created, nil
}

func (s *WishlistService) GetByID(ctx context.Context, userID, wishlistID int64) (domain.Wishlist, error) {
	w, err := s.wishlistRepo.FindByID(ctx, wishlistID)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("find wishlist: %w", err)
	}

	if w.UserID != userID {
		return domain.Wishlist{}, domain.ErrWishlistForbidden
	}
	return w, nil
}

func (s *WishlistService) GetAllByUser(ctx context.Context, userID int64) ([]domain.Wishlist, error) {
	wishlists, err := s.wishlistRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find wishlists by user id: %w", err)
	}
	return wishlists, nil
}

func (s *WishlistService) Update(ctx context.Context, w domain.Wishlist) (domain.Wishlist, error) {
	existing, err := s.wishlistRepo.FindByID(ctx, w.ID)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("find wishlist for update: %w", err)
	}

	if existing.UserID != w.UserID {
		return domain.Wishlist{}, domain.ErrWishlistForbidden
	}

	w.PublicToken = existing.PublicToken

	updated, err := s.wishlistRepo.Update(ctx, w)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("update wishlist: %w", err)
	}
	return updated, nil
}

func (s *WishlistService) Delete(ctx context.Context, userID, wishlistID int64) error {
	existing, err := s.wishlistRepo.FindByID(ctx, wishlistID)
	if err != nil {
		return fmt.Errorf("find wishlist for delete: %w", err)
	}

	if existing.UserID != userID {
		return domain.ErrWishlistForbidden
	}

	if err := s.wishlistRepo.Delete(ctx, wishlistID); err != nil {
		return fmt.Errorf("delete wishlist: %w", err)
	}
	return nil
}

func (s *WishlistService) GetByPublicToken(ctx context.Context, token string) (domain.Wishlist, error) {
	wishlist, err := s.wishlistRepo.FindByPublicToken(ctx, token)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("find wishlist by public token: %w", err)
	}
	return wishlist, nil
}
