package repositories

import (
	"blog-api/Domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type passwordResetTokenModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	TokenHash string             `bson:"token_hash"`
	ExpiresAt time.Time          `bson:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"`
	Used      bool               `bson:"used"`
}

type passwordResetTokenRepo struct {
	collection *mongo.Collection
}

func NewPasswordResetTokenRepo(db *mongo.Database) domain.IPasswordResetTokenRepository {
	return &passwordResetTokenRepo{
		collection: db.Collection("password_reset_tokens"),
	}
}

func (r *passwordResetTokenRepo) Store(ctx context.Context, token *domain.PasswordResetToken) error {
	doc := passwordResetTokenModel{
		ID:        primitive.NewObjectID(),
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
		Used:      false,
	}
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *passwordResetTokenRepo) GetByTokenHash(ctx context.Context, rawToken string) (*domain.PasswordResetToken, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc passwordResetTokenModel
		if err := cursor.Decode(&doc); err != nil {
			continue 
		}

		if err := bcrypt.CompareHashAndPassword([]byte(doc.TokenHash), []byte(rawToken)); err == nil {
			return &domain.PasswordResetToken{
				ID:        doc.ID.Hex(),
				UserID:    doc.UserID,
				TokenHash: doc.TokenHash,
				ExpiresAt: doc.ExpiresAt,
				CreatedAt: doc.CreatedAt,
				Used:      doc.Used,
			}, nil
		}
	}

	return nil, fmt.Errorf("reset token not found or invalid")
}


func (r *passwordResetTokenRepo) MarkUsed(ctx context.Context, id string) error {
	objID, err:= primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"used": true}})
	return err
}

func (r *passwordResetTokenRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": time.Now()}})
	return err
}

