package infrastructure

import (
	domain "blog-api/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTokenRepository struct {
	Collection *mongo.Collection
}

func NewMongoTokenRepository(db *mongo.Database) *MongoTokenRepository {
	return &MongoTokenRepository{
		Collection: db.Collection("tokens"),
	}
}

func (r *MongoTokenRepository) Save(ctx context.Context, token *domain.Token) error {
	_, err := r.Collection.InsertOne(ctx, token)
	return err
}

func (r *MongoTokenRepository) FindByToken(ctx context.Context, rawToken string) (*domain.Token, error) {
	var token domain.Token
	err := r.Collection.FindOne(ctx, bson.M{
		"token":   rawToken,
		"revoked": false,
	}).Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *MongoTokenRepository) Revoke(ctx context.Context, rawToken string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"token": rawToken}, bson.M{"$set": bson.M{"revoked": true}})
	return err
}
