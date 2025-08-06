package repositories

import (
	domain "blog-api/Domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *blogRepository) FindByID(ctx context.Context, blogID string) (*domain.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID format: %w", err)
	}

	var blogDoc blogModel
	err = r.blogCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&blogDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("blog not found")
		}
		return nil, fmt.Errorf("failed to find blog: %w", err)
	}

	domainBlog := toDomainBlog(blogDoc)
	return &domainBlog, nil
}

func (r *blogRepository) Update(ctx context.Context, blog *domain.Blog) (*domain.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(blog.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"title":     blog.Title,
			"content":   blog.Content,
			"tags":      blog.Tags,
			"updatedAt": blog.UpdatedAt,
		},
	}

	_, err = r.blogCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}

	return blog, nil
}

func (r *blogRepository) Filter(ctx context.Context, tag string, date string, sort string) ([]domain.Blog, error) {
	filter := bson.M{}

	if tag != "" {
		filter["tags"] = tag
	}

	if date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err == nil {
			start := parsedDate
			end := parsedDate.Add(24 * time.Hour)
			filter["createdAt"] = bson.M{
				"$gte": start,
				"$lt":  end,
			}
		}
	}

	findOptions := options.Find()
	if sort == "popular" {
		findOptions.SetSort(bson.D{{"view_count", -1}})
	} else {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cursor, err := r.blogCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blogs: %w", err)
	}
	defer cursor.Close(ctx)

	var results []blogModel
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode blogs: %w", err)
	}

	var blogs []domain.Blog
	for _, b := range results {
		blogs = append(blogs, toDomainBlog(b))
	}
	return blogs, nil
}
