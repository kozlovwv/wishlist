package application

import (
	"context"

	"wishlist/internal/domain"
)

type ItemService struct {
	itemRepo     ItemRepository
	wishlistRepo WishlistRepository
}

func NewItemService(itemRepo ItemRepository, wishlistRepo WishlistRepository) *ItemService {
	return &ItemService{
		itemRepo:     itemRepo,
		wishlistRepo: wishlistRepo,
	}
}

func (s *ItemService) Create(ctx context.Context, userID int64, item domain.Item) (domain.Item, error) {
	w, err := s.wishlistRepo.FindByID(ctx, item.WishlistID)
	if err != nil {
		return domain.Item{}, err
	}

	if w.UserID != userID {
		return domain.Item{}, domain.ErrWishlistForbidden
	}

	return s.itemRepo.Create(ctx, item)
}

func (s *ItemService) GetAllByWishlist(ctx context.Context, userID, wishlistID int64) ([]domain.Item, error) {
	w, err := s.wishlistRepo.FindByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	if w.UserID != userID {
		return nil, domain.ErrWishlistForbidden
	}

	return s.itemRepo.FindByWishlistID(ctx, wishlistID)
}

func (s *ItemService) Update(ctx context.Context, userID int64, item domain.Item) (domain.Item, error) {
	existing, err := s.itemRepo.FindByID(ctx, item.ID)
	if err != nil {
		return domain.Item{}, err
	}

	w, err := s.wishlistRepo.FindByID(ctx, existing.WishlistID)
	if err != nil {
		return domain.Item{}, err
	}

	if w.UserID != userID {
		return domain.Item{}, domain.ErrWishlistForbidden
	}

	item.WishlistID = existing.WishlistID
	return s.itemRepo.Update(ctx, item)
}

func (s *ItemService) Delete(ctx context.Context, userID, itemID int64) error {
	existing, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	w, err := s.wishlistRepo.FindByID(ctx, existing.WishlistID)
	if err != nil {
		return err
	}

	if w.UserID != userID {
		return domain.ErrWishlistForbidden
	}

	return s.itemRepo.Delete(ctx, itemID)
}

func (s *ItemService) Reserve(ctx context.Context, publicToken string, itemID int64) error {
	w, err := s.wishlistRepo.FindByPublicToken(ctx, publicToken)
	if err != nil {
		return err
	}

	item, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	if item.WishlistID != w.ID {
		return domain.ErrItemNotFound
	}

	return s.itemRepo.Reserve(ctx, itemID)
}
