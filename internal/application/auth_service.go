package application

import (
	"context"
	"errors"
	"fmt"

	"wishlist/internal/domain"
)

type AuthService struct {
	userRepo     UserRepository
	hasher       PasswordHasher
	tokenManager TokenManager
}

func NewAuthService(userRepo UserRepository, hasher PasswordHasher, tokenManager TokenManager) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return "", domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return "", fmt.Errorf("find user by email: %w", err)
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, domain.User{
		Email:        email,
		PasswordHash: hash,
	})
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	token, err := s.tokenManager.Generate(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return "", domain.ErrUserNotFound
	}
	if err != nil {
		return "", fmt.Errorf("find user by email: %w", err)
	}

	if !s.hasher.Verify(password, user.PasswordHash) {
		return "", domain.ErrInvalidPassword
	}

	token, err := s.tokenManager.Generate(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}
