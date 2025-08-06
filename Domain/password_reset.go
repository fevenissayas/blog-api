package domain

import (
	"context"
	"time"
)

type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	Used      bool
}

type IPasswordResetTokenRepository interface {
	Store(ctx context.Context, token *PasswordResetToken) error
	GetByTokenHash(ctx context.Context, hash string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

// type IPasswordResetUsecase interface {
// 	RequestPasswordReset(ctx context.Context, input RequestPasswordResetInput) (token string, err error)
// }
