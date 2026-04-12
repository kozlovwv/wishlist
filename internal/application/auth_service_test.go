package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"wishlist/internal/application"
	"wishlist/internal/domain"
	mock_application "wishlist/mock"
)

func TestAuthService_Register(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_application.MockUserRepository,
		hasher *mock_application.MockPasswordHasher,
		tokenManager *mock_application.MockTokenManager,
	)

	tests := []struct {
		name          string
		email         string
		password      string
		mockBehavior  mockBehavior
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "success",
			email:    "test@mail.com",
			password: "password",
			mockBehavior: func(userRepo *mock_application.MockUserRepository, hasher *mock_application.MockPasswordHasher, tokenManager *mock_application.MockTokenManager) {
				ctx := context.Background()

				userRepo.EXPECT().
					FindByEmail(ctx, "test@mail.com").
					Return(domain.User{}, domain.ErrUserNotFound)

				hasher.EXPECT().
					Hash("password").
					Return("hashed", nil)

				userRepo.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, u domain.User) (domain.User, error) {
						u.ID = 1
						return u, nil
					})

				tokenManager.EXPECT().
					Generate(int64(1)).
					Return("token-123", nil)
			},
			expectedToken: "token-123",
			expectedErr:   nil,
		},
		{
			name:     "user already exists",
			email:    "test@mail.com",
			password: "password",
			mockBehavior: func(userRepo *mock_application.MockUserRepository, hasher *mock_application.MockPasswordHasher, tokenManager *mock_application.MockTokenManager) {
				ctx := context.Background()

				userRepo.EXPECT().
					FindByEmail(ctx, "test@mail.com").
					Return(domain.User{ID: 1}, nil)
			},
			expectedToken: "",
			expectedErr:   domain.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_application.NewMockUserRepository(ctrl)
			hasher := mock_application.NewMockPasswordHasher(ctrl)
			tokenManager := mock_application.NewMockTokenManager(ctrl)

			tt.mockBehavior(userRepo, hasher, tokenManager)

			svc := application.NewAuthService(userRepo, hasher, tokenManager)

			token, err := svc.Register(context.Background(), tt.email, tt.password)

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_application.MockUserRepository,
		hasher *mock_application.MockPasswordHasher,
		tokenManager *mock_application.MockTokenManager,
	)

	tests := []struct {
		name          string
		email         string
		password      string
		mockBehavior  mockBehavior
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "success",
			email:    "test@mail.com",
			password: "password",
			mockBehavior: func(userRepo *mock_application.MockUserRepository, hasher *mock_application.MockPasswordHasher, tokenManager *mock_application.MockTokenManager) {
				userRepo.EXPECT().
					FindByEmail(gomock.Any(), "test@mail.com").
					Return(domain.User{
						ID:           1,
						PasswordHash: "hash",
					}, nil)

				hasher.EXPECT().
					Verify("password", "hash").
					Return(true)

				tokenManager.EXPECT().
					Generate(int64(1)).
					Return("jwt-token", nil)
			},
			expectedToken: "jwt-token",
			expectedErr:   nil,
		},
		{
			name:     "wrong password",
			email:    "test@mail.com",
			password: "wrong",
			mockBehavior: func(userRepo *mock_application.MockUserRepository, hasher *mock_application.MockPasswordHasher, tokenManager *mock_application.MockTokenManager) {
				userRepo.EXPECT().
					FindByEmail(gomock.Any(), "test@mail.com").
					Return(domain.User{
						ID:           1,
						PasswordHash: "hash",
					}, nil)

				hasher.EXPECT().
					Verify("wrong", "hash").
					Return(false)
			},
			expectedToken: "",
			expectedErr:   domain.ErrInvalidPassword,
		},
		{
			name:     "user not found",
			email:    "test@mail.com",
			password: "password",
			mockBehavior: func(userRepo *mock_application.MockUserRepository, hasher *mock_application.MockPasswordHasher, tokenManager *mock_application.MockTokenManager) {
				userRepo.EXPECT().
					FindByEmail(gomock.Any(), "test@mail.com").
					Return(domain.User{}, domain.ErrUserNotFound)
			},
			expectedToken: "",
			expectedErr:   domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_application.NewMockUserRepository(ctrl)
			hasher := mock_application.NewMockPasswordHasher(ctrl)
			tokenManager := mock_application.NewMockTokenManager(ctrl)

			tt.mockBehavior(userRepo, hasher, tokenManager)

			svc := application.NewAuthService(userRepo, hasher, tokenManager)

			token, err := svc.Login(context.Background(), tt.email, tt.password)

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedToken, token)
		})
	}
}
