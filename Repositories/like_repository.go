package repositories

import (
	"blog-api/Domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LikeRepository struct {
	collection *mongo.Collection
}

func NewLikeRepository(db *mongo.Database) domain.IBlogLikeRepository {
	collection := db.Collection("likes")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "blogId", Value: 1},
			{Key: "userId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	collection.Indexes().CreateOne(context.Background(), indexModel)

	return &LikeRepository{collection}
}

func (r *LikeRepository) AddLike(like *domain.Like) error {
	blogOID, err := primitive.ObjectIDFromHex(like.BlogID)
	if err != nil {
		return domain.ErrInvalidInput
	}
	userOID, err := primitive.ObjectIDFromHex(like.UserID)
	if err != nil {
		return domain.ErrInvalidInput
	}

	doc := bson.M{
		"blogId":    blogOID,
		"userId":    userOID,
		"createdAt": time.Now(),
	}
	_, err = r.collection.InsertOne(context.Background(), doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (r *LikeRepository) RemoveLike(blogID, userID string) error {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return domain.ErrInvalidInput
	}
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrInvalidInput
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{
		"blogId": blogOID,
		"userId": userOID,
	})
	return err
}

func (r *LikeRepository) IsLiked(blogID, userID string) (bool, error) {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, domain.ErrInvalidInput
	}
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, domain.ErrInvalidInput
	}

	count, err := r.collection.CountDocuments(context.Background(), bson.M{
		"blogId": blogOID,
		"userId": userOID,
	})
	return count > 0, err
}

func (r *LikeRepository) CountLikes(blogID string) (int, error) {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return 0, domain.ErrInvalidInput
	}

	count, err := r.collection.CountDocuments(context.Background(), bson.M{"blogId": blogOID})
	return int(count), err
}
