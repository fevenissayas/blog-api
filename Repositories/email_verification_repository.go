package repositories

import (
	domain "blog-api/Domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type emailVerificationModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	OTP       string             `bson:"otp"`
	ExpiresAt time.Time          `bson:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"`
	Used      bool               `bson:"used"`
}

type emailVerificationRepo struct {
	collection *mongo.Collection
}

func NewEmailVerificationRepository(db *mongo.Database) domain.IEmailVerificationRepository {
	collection := db.Collection("email_verifications")
	
	// Create TTL index to auto-delete expired documents
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "expires_at", Value: 1}},
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)
	
	return &emailVerificationRepo{
		collection: collection,
	}
}

func (r *emailVerificationRepo) Store(ctx context.Context, verification *domain.EmailVerification) error {
	objID := primitive.NewObjectID()
	verification.ID = objID.Hex()

	doc := emailVerificationModel{
		ID:        objID,
		Email:     verification.Email,
		OTP:       verification.OTP,
		ExpiresAt: verification.ExpiresAt,
		CreatedAt: verification.CreatedAt,
		Used:      false,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to store email verification: %w", err)
	}

	return nil
}

func (r *emailVerificationRepo) GetByEmail(ctx context.Context, email string) (*domain.EmailVerification, error) {
	var doc emailVerificationModel
	
	filter := bson.M{
		"email":      email,
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no valid verification found for email")
		}
		return nil, fmt.Errorf("failed to find email verification: %w", err)
	}

	return &domain.EmailVerification{
		ID:        doc.ID.Hex(),
		Email:     doc.Email,
		OTP:       doc.OTP,
		ExpiresAt: doc.ExpiresAt,
		CreatedAt: doc.CreatedAt,
		Used:      doc.Used,
	}, nil
}

func (r *emailVerificationRepo) GetByOTP(ctx context.Context, otp, email string) (*domain.EmailVerification, error) {
	var doc emailVerificationModel
	
	filter := bson.M{
		"email":      email,
		"otp":        otp,
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid or expired OTP")
		}
		return nil, fmt.Errorf("failed to find email verification: %w", err)
	}

	return &domain.EmailVerification{
		ID:        doc.ID.Hex(),
		Email:     doc.Email,
		OTP:       doc.OTP,
		ExpiresAt: doc.ExpiresAt,
		CreatedAt: doc.CreatedAt,
		Used:      doc.Used,
	}, nil
}

func (r *emailVerificationRepo) MarkUsed(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid verification ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{"used": true},
	}

	_, err = r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return fmt.Errorf("failed to mark verification as used: %w", err)
	}

	return nil
}

func (r *emailVerificationRepo) DeleteExpired(ctx context.Context) error {
	filter := bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete expired verifications: %w", err)
	}

	return nil
}
