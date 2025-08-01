package domain

import (
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