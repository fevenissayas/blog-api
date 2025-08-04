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

type blogModel struct {
	ID             primitive.ObjectID `bson:"_id"`
	Title          string             `bson:"title"`
	Content        string             `bson:"content"`
	UserID         string		  	  `bson:"user_id"`
	Tags           []string			  `bson:"tags"`
	ViewCount 	   int				  `bson:"view_count"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

func toDomainBlog(m blogModel) domain.Blog {
	return domain.Blog{
		ID:             m.ID.Hex(),
		Title:          m.Title,
		Content:        m.Content,
		Tags:			m.Tags,
		UserID:			m.UserID,
		ViewCount: 		m.ViewCount,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

type blogRepository struct {
	blogCollection *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) domain.IBlogRepository {

	collection := db.Collection("blogs")
	return &blogRepository{blogCollection: collection}
}

func (r *blogRepository) Create(ctx context.Context, blog *domain.Blog) (*domain.Blog, error) {
	blogObjectID := primitive.NewObjectID()
	blog.ID = blogObjectID.Hex()

	doc := bson.M{
		"_id":             blogObjectID,
		"title":           blog.Title,
		"content":         blog.Content,
		"tags":            blog.Tags,
		"view_count": 	   blog.ViewCount,
		"user_id":         blog.UserID,
		"createdAt":       blog.CreatedAt,
		"updatedAt":       blog.UpdatedAt,
	}

	_, err := r.blogCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert blog: %w", err)
	}

	return blog, nil
}
