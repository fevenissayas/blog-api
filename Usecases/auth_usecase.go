package usecases

import (
	domain "blog-api/Domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type AuthUsecase struct {
	jwtService       domain.IJWTService
	userRepo         domain.IUserRepository
	refreshTokenRepo domain.IRefreshTokenRepository
	contextTimeout   time.Duration
}

func NewAuthUsecase(
	jwtService domain.IJWTService,
	userRepo domain.IUserRepository,
	refreshTokenRepo domain.IRefreshTokenRepository,
	timeout time.Duration,
) domain.IAuthUsecase {
	return &AuthUsecase{
		jwtService:       jwtService,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		contextTimeout:   timeout,
	}
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, token string) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	payload, err := uc.jwtService.ValidateRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	storedToken, err := uc.refreshTokenRepo.FindByID(ctx, payload.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch refresh token: %w", err)
	}
	if storedToken == nil || storedToken.Token != token {
		return nil, errors.New("refresh token not found or mismatched")
	}

	if storedToken.RevokedAt != nil {
		return nil, errors.New("refresh token has been revoked")
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	user, err := uc.userRepo.GetByID(ctx, payload.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	newAccessToken, err := uc.jwtService.CreateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	newRefreshToken, tokenPayload, err := uc.jwtService.CreateRefreshToken(user)
if err != nil {
    return nil, fmt.Errorf("failed to create refresh token: %w", err)
}

newRT := &domain.RefreshToken{
    ID:        tokenPayload.TokenID,
    Token:     newRefreshToken,
    UserID:    user.ID,
    CreatedAt: tokenPayload.IssuedAt,
    UpdatedAt: tokenPayload.IssuedAt,
    ExpiresAt: tokenPayload.ExpiresAt,
}

	if err := uc.refreshTokenRepo.StoreToken(ctx, newRT); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	if err := uc.refreshTokenRepo.RevokeToken(ctx, storedToken.Token); err != nil {
		log.Printf("warning: failed to revoke old refresh token: %v", err)
	}

	return &domain.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
