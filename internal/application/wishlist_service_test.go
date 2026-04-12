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

func TestWishlistService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wishlistRepo := mock_application.NewMockWishlistRepository(ctrl)
	tokenGen := mock_application.NewMockPublicTokenGenerator(ctrl)

	svc := application.NewWishlistService(wishlistRepo, tokenGen)

	input := domain.Wishlist{
		UserID: 1,
		Title:  "test",
	}

	tokenGen.EXPECT().
		Generate().
		Return("token-123", nil)

	wishlistRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, w domain.Wishlist) (domain.Wishlist, error) {
			require.Equal(t, "token-123", w.PublicToken)
			return domain.Wishlist{
				ID:          1,
				UserID:      1,
				Title:       "test",
				PublicToken: "token-123",
			}, nil
		})

	res, err := svc.Create(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, int64(1), res.ID)
}

func TestWishlistService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		setup     func(repo *mock_application.MockWishlistRepository)
		expectErr error
	}{
		{
			name:   "success",
			userID: 1,
			setup: func(repo *mock_application.MockWishlistRepository) {
				repo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)
			},
		},
		{
			name:   "forbidden",
			userID: 2,
			setup: func(repo *mock_application.MockWishlistRepository) {
				repo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{ID: 10, UserID: 1}, nil)
			},
			expectErr: domain.ErrWishlistForbidden,
		},
		{
			name:   "repo error wrapped",
			userID: 1,
			setup: func(repo *mock_application.MockWishlistRepository) {
				repo.EXPECT().
					FindByID(gomock.Any(), int64(10)).
					Return(domain.Wishlist{}, errors.New("db error"))
			},
			expectErr: errors.New("find wishlist: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_application.NewMockWishlistRepository(ctrl)
			svc := application.NewWishlistService(repo, nil)

			tt.setup(repo)

			_, err := svc.GetByID(context.Background(), tt.userID, 10)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWishlistService_GetAllByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_application.NewMockWishlistRepository(ctrl)
	svc := application.NewWishlistService(repo, nil)

	repo.EXPECT().
		FindByUserID(gomock.Any(), int64(1)).
		Return([]domain.Wishlist{
			{ID: 1, UserID: 1},
			{ID: 2, UserID: 1},
		}, nil)

	res, err := svc.GetAllByUser(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, res, 2)
}

func TestWishlistService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_application.NewMockWishlistRepository(ctrl)
	svc := application.NewWishlistService(repo, nil)

	input := domain.Wishlist{
		ID:     10,
		UserID: 1,
		Title:  "new",
	}

	repo.EXPECT().
		FindByID(gomock.Any(), int64(10)).
		Return(domain.Wishlist{
			ID:          10,
			UserID:      1,
			PublicToken: "token",
		}, nil)

	repo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(domain.Wishlist{
			ID:     10,
			UserID: 1,
			Title:  "new",
		}, nil)

	res, err := svc.Update(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, int64(10), res.ID)
}

func TestWishlistService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_application.NewMockWishlistRepository(ctrl)
	svc := application.NewWishlistService(repo, nil)

	repo.EXPECT().
		FindByID(gomock.Any(), int64(10)).
		Return(domain.Wishlist{ID: 10, UserID: 1}, nil)

	repo.EXPECT().
		Delete(gomock.Any(), int64(10)).
		Return(nil)

	err := svc.Delete(context.Background(), 1, 10)

	require.NoError(t, err)
}

func TestWishlistService_GetByPublicToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_application.NewMockWishlistRepository(ctrl)
	svc := application.NewWishlistService(repo, nil)

	repo.EXPECT().
		FindByPublicToken(gomock.Any(), "token").
		Return(domain.Wishlist{ID: 1}, nil)

	res, err := svc.GetByPublicToken(context.Background(), "token")

	require.NoError(t, err)
	require.Equal(t, int64(1), res.ID)
}
