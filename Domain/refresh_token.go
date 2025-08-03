package domain

import (
	"context"
	"time"
)


type RefreshToken struct{
	ID 			string
	Token 		string
	UserID		string
	CreatedAt 	time.Time
	UpdatedAt	time.Time
	ExpiresAt	time.Time
	RevokedAt   *time.Time
}

type IAuthUsecase interface {
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

type IRefreshTokenRepository interface {
	StoreToken(ctx context.Context, token *RefreshToken) error
	FindByID(ctx context.Context, tokenID string) (*RefreshToken, error)
	RevokeToken(ctx context.Context, token string) error
	DeleteAllTokensForUser(ctx context.Context, userID string) error
	// DeleteByUserID(ctx context.Context, userID string) error //i'll implement it later when i do user update
}

