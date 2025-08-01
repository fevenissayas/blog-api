package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Token struct {
	ID        primitive.ObjectID
	UserID    primitive.ObjectID
	Token     string
	Revoked   bool
	CreatedAt time.Time
	ExpiresAt time.Time
}

type TokenRepository interface {
	Save(ctx context.Context, token *Token) error
	FindByToken(ctx context.Context, token string) (*Token, error)
	Revoke(ctx context.Context, token string) error
}
