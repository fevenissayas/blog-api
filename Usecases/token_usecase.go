package usecases

import (
	domain "blog-api/Domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenUsecase struct {
	Repo domain.TokenRepository
}

func NewTokenUsecase(repo domain.TokenRepository) *TokenUsecase {
	return &TokenUsecase{Repo: repo}
}

func (u *TokenUsecase) SaveRefreshToken(ctx context.Context, userID primitive.ObjectID, token string, expiresAt time.Time) error {
	return u.Repo.Save(ctx, &domain.Token{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	})
}

func (u *TokenUsecase) FindValidRefreshToken(ctx context.Context, token string) (*domain.Token, error) {
	return u.Repo.FindByToken(ctx, token)
}

func (u *TokenUsecase) RevokeToken(ctx context.Context, token string) error {
	return u.Repo.Revoke(ctx, token)
}
