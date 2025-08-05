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

type userModel struct {
	ID             primitive.ObjectID `bson:"_id"`
	Email          string             `bson:"email"`
	Username       string             `bson:"username"`
	Password       string             `bson:"password"`
	Role           domain.Role        `bson:"role"`
	IsVerified     bool               `bson:"isVerified"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	Bio            string             `bson:"bio,omitempty"`
	ProfilePicture string             `bson:"profile_picture,omitempty"`
	ContactInfo    string             `bson:"contact_info,omitempty"`
}

func toDomainUser(m userModel) domain.User {
	return domain.User{
		ID:             m.ID.Hex(),
		Email:          m.Email,
		Username:       m.Username,
		Password:       m.Password,
		Role:           m.Role,
		IsVerified:     m.IsVerified,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		Bio:            m.Bio,
		ProfilePicture: m.ProfilePicture,
		ContactInfo:    m.ContactInfo,
	}
}

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.IUserRepository {

	collection := db.Collection("user")
	return &userRepository{userCollection: collection}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	userObjectID := primitive.NewObjectID()
	user.ID = userObjectID.Hex()

	doc := bson.M{
		"_id":             userObjectID,
		"email":           user.Email,
		"username":        user.Username,
		"password":        user.Password,
		"role":            user.Role,
		"isVerified":      user.IsVerified,
		"createdAt":       user.CreatedAt,
		"updatedAt":       user.UpdatedAt,
		"bio":             user.Bio,
		"profile_picture": user.ProfilePicture,
		"contact_info":    user.ContactInfo,
	}

	_, err := r.userCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func (r *userRepository) findUserByField(ctx context.Context, field string, value string) (*domain.User, error) {
	var model userModel
	filter := bson.M{field: value}
	err := r.userCollection.FindOne(ctx, filter).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user by %s: %w", field, err)
	}
	user := toDomainUser(model)
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findUserByField(ctx, "email", email)
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.findUserByField(ctx, "username", username)
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var model userModel
	err = r.userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	user := toDomainUser(model)
	return &user, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, fmt.Errorf("failed to count documents by email: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.userCollection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, fmt.Errorf("failed to count documents by username: %w", err)
	}
	return count > 0, nil
}
func (r *userRepository) Promote (ctx context.Context, user domain.User)(error){
	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %w", err)
	}
	update := bson.M{
		"$set": bson.M{
			"role": domain.RoleAdmin,
			"updatedAt": user.UpdatedAt,
		},
	}
	_, err = r.userCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}
	return nil
}
