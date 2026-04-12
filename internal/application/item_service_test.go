package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"wishlist/internal/application"
	"wishlist/internal/domain"
	mock_application "wishlist/mock"
)

func TestItemService_Create(t *testing.T) {
	type mockBehavior func(
		wishlistRepo *mock_application.MockWishlistRepository,
		itemRepo *mock_application.MockItemRepository,
	)

	tests := []struct {
		name      string
		userID    int64
		item      domain.Item
		mock      mockBehavior
		expectErr error
	}{
		{
			name:   "success",
			userID: 1,
			item: domain.Item{
				WishlistID: 10,
				Title:      "item",
			},
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)

				itemRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(domain.Item{ID: 1, WishlistID: 10, Title: "item"}, nil)
			},
			expectErr: nil,
		},
		{
			name:   "forbidden wishlist",
			userID: 2,
			item: domain.Item{
				WishlistID: 10,
			},
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)
			},
			expectErr: domain.ErrWishlistForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wishlistRepo := mock_application.NewMockWishlistRepository(ctrl)
			itemRepo := mock_application.NewMockItemRepository(ctrl)

			tt.mock(wishlistRepo, itemRepo)

			svc := application.NewItemService(itemRepo, wishlistRepo)

			_, err := svc.Create(context.Background(), tt.userID, tt.item)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestItemService_GetAllByWishlist(t *testing.T) {
	type mockBehavior func(
		wishlistRepo *mock_application.MockWishlistRepository,
		itemRepo *mock_application.MockItemRepository,
	)

	tests := []struct {
		name       string
		userID     int64
		wishlistID int64
		mock       mockBehavior
		expectErr  error
	}{
		{
			name:       "success",
			userID:     1,
			wishlistID: 10,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)

				itemRepo.EXPECT().
					FindByWishlistID(gomock.Any(), int64(10)).
					Return([]domain.Item{
						{ID: 1, WishlistID: 10},
					}, nil)
			},
			expectErr: nil,
		},
		{
			name:       "forbidden wishlist",
			userID:     2,
			wishlistID: 10,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)
			},
			expectErr: domain.ErrWishlistForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wishlistRepo := mock_application.NewMockWishlistRepository(ctrl)
			itemRepo := mock_application.NewMockItemRepository(ctrl)

			tt.mock(wishlistRepo, itemRepo)

			svc := application.NewItemService(itemRepo, wishlistRepo)

			_, err := svc.GetAllByWishlist(context.Background(), tt.userID, tt.wishlistID)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestItemService_Update(t *testing.T) {
	type mockBehavior func(
		wishlistRepo *mock_application.MockWishlistRepository,
		itemRepo *mock_application.MockItemRepository,
	)

	tests := []struct {
		name      string
		userID    int64
		item      domain.Item
		mock      mockBehavior
		expectErr error
	}{
		{
			name:   "success",
			userID: 1,
			item: domain.Item{
				ID:    1,
				Title: "updated",
			},
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				itemRepo.EXPECT().
					FindByID(gomock.Any(), int64(1)).
					Return(domain.Item{ID: 1, WishlistID: 10}, nil)

				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)

				itemRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(domain.Item{ID: 1, WishlistID: 10, Title: "updated"}, nil)
			},
			expectErr: nil,
		},
		{
			name:   "forbidden",
			userID: 2,
			item: domain.Item{
				ID: 1,
			},
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				itemRepo.EXPECT().
					FindByID(gomock.Any(), int64(1)).
					Return(domain.Item{ID: 1, WishlistID: 10}, nil)

				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)
			},
			expectErr: domain.ErrWishlistForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wishlistRepo := mock_application.NewMockWishlistRepository(ctrl)
			itemRepo := mock_application.NewMockItemRepository(ctrl)

			tt.mock(wishlistRepo, itemRepo)

			svc := application.NewItemService(itemRepo, wishlistRepo)

			_, err := svc.Update(context.Background(), tt.userID, tt.item)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestItemService_Delete(t *testing.T) {
	type mockBehavior func(
		wishlistRepo *mock_application.MockWishlistRepository,
		itemRepo *mock_application.MockItemRepository,
	)

	tests := []struct {
		name      string
		userID    int64
		itemID    int64
		mock      mockBehavior
		expectErr error
	}{
		{
			name:   "success",
			userID: 1,
			itemID: 1,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				itemRepo.EXPECT().
					FindByID(gomock.Any(), int64(1)).
					Return(domain.Item{ID: 1, WishlistID: 10}, nil)

				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)

				itemRepo.EXPECT().
					Delete(gomock.Any(), int64(1)).
					Return(nil)
			},
			expectErr: nil,
		},
		{
			name:   "forbidden",
			userID: 2,
			itemID: 1,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				itemRepo.EXPECT().
					FindByID(gomock.Any(), int64(1)).
					Return(domain.Item{ID: 1, WishlistID: 10}, nil)

				wishlistRepo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)
			},
			expectErr: domain.ErrWishlistForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wishlistRepo := mock_application.NewMockWishlistRepository(ctrl)
			itemRepo := mock_application.NewMockItemRepository(ctrl)

			tt.mock(wishlistRepo, itemRepo)

			svc := application.NewItemService(itemRepo, wishlistRepo)

			err := svc.Delete(context.Background(), tt.userID, tt.itemID)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestItemService_Reserve(t *testing.T) {
	type mockBehavior func(
		wishlistRepo *mock_application.MockWishlistRepository,
		itemRepo *mock_application.MockItemRepository,
	)

	tests := []struct {
		name      string
		token     string
		itemID    int64
		mock      mockBehavior
		expectErr error
	}{
		{
			name:   "success reserve",
			token:  "public-token",
			itemID: 1,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByPublicToken(gomock.Any(), "public-token").
					Return(domain.Wishlist{ID: 10}, nil)

				itemRepo.EXPECT().
					FindByID(gomock.Any(), int64(1)).
					Return(domain.Item{ID: 1, WishlistID: 10}, nil)

				itemRepo.EXPECT().
					Reserve(gomock.Any(), int64(1)).
					Return(nil)
			},
			expectErr: nil,
		},
		{
			name:   "item not in wishlist",
			token:  "public-token",
			itemID: 1,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByPublicToken(gomock.Any(), "public-token").
					Return(domain.Wishlist{ID: 10}, nil)

				itemRepo.EXPECT().
					FindByID(gomock.Any(), int64(1)).
					Return(domain.Item{ID: 1, WishlistID: 99}, nil)
			},
			expectErr: domain.ErrItemNotFound,
		},
		{
			name:   "wishlist not found",
			token:  "bad-token",
			itemID: 1,
			mock: func(wishlistRepo *mock_application.MockWishlistRepository, itemRepo *mock_application.MockItemRepository) {
				wishlistRepo.EXPECT().
					FindByPublicToken(gomock.Any(), "bad-token").
					Return(domain.Wishlist{}, errors.New("not found"))
			},
			expectErr: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wishlistRepo := mock_application.NewMockWishlistRepository(ctrl)
			itemRepo := mock_application.NewMockItemRepository(ctrl)

			tt.mock(wishlistRepo, itemRepo)

			svc := application.NewItemService(itemRepo, wishlistRepo)

			err := svc.Reserve(context.Background(), tt.token, tt.itemID)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
