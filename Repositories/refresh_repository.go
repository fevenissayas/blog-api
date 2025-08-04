package repositories

import (
	"context"
	"fmt"
	"time"

	domain "blog-api/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type refreshTokenModel struct {
    ID        string`bson:"_id,omitempty"`
    Token     string             `bson:"token"`
    UserID    string             `bson:"userId"`
    CreatedAt time.Time          `bson:"createdAt"`
    UpdatedAt time.Time          `bson:"updatedAt"`
    ExpiresAt time.Time          `bson:"expiresAt"`
    RevokedAt *time.Time         `bson:"revokedAt,omitempty"`
}

func toDomainRefreshToken(m *refreshTokenModel) *domain.RefreshToken {
    return &domain.RefreshToken{
        ID:        m.ID,
        Token:     m.Token,
        UserID:    m.UserID,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
        ExpiresAt: m.ExpiresAt,
        RevokedAt: m.RevokedAt,
    }
}

type refreshTokenRepository struct {
    refreshTokenCollection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) domain.IRefreshTokenRepository {
   collection := db.Collection("refresh_tokens")

    return &refreshTokenRepository{refreshTokenCollection: collection}
}

func (r *refreshTokenRepository) StoreToken(ctx context.Context, token *domain.RefreshToken) error {
    doc := bson.M{
        "_id":       token.ID, 
        "token":     token.Token,
        "userId":    token.UserID,
        "createdAt": token.CreatedAt,
        "updatedAt": token.UpdatedAt,
        "expiresAt": token.ExpiresAt,
    }

    if token.RevokedAt != nil {
        doc["revokedAt"] = token.RevokedAt
    }

    _, err := r.refreshTokenCollection.InsertOne(ctx, doc)
    if err != nil {
        return fmt.Errorf("failed to insert refresh token: %w", err)
    }

    return nil
}


func (r *refreshTokenRepository) FindByID(ctx context.Context, tokenID string) (*domain.RefreshToken, error) {
     filter := bson.M{"_id": tokenID} 
    var model refreshTokenModel
    err := r.refreshTokenCollection.FindOne(ctx, filter).Decode(&model)
    if err == mongo.ErrNoDocuments {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find refresh token by ID: %w", err)
    }
    return toDomainRefreshToken(&model), nil
}

func (r *refreshTokenRepository) DeleteAllTokensForUser(ctx context.Context, userID string) error {
	filter := bson.M{"userId": userID}

	_, err := r.refreshTokenCollection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete refresh tokens for user %s: %w", userID, err)
	}
	return nil
}


func (r *refreshTokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
    now := time.Now()
    filter := bson.M{"token": tokenString}
    update := bson.M{
        "$set": bson.M{
            "revokedAt": now,
            "updatedAt": now,
        },
    }

    result, err := r.refreshTokenCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to revoke token: %w", err)
    }
    if result.MatchedCount == 0 {
        return fmt.Errorf("refresh token not found")
    }
    return nil
}
