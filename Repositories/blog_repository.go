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
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	UserID    string             `bson:"user_id"`
	Tags      []string           `bson:"tags"`
	ViewCount int                `bson:"view_count"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func toDomainBlog(m blogModel) domain.Blog {
	return domain.Blog{
		ID:        m.ID.Hex(),
		Title:     m.Title,
		Content:   m.Content,
		Tags:      m.Tags,
		UserID:    m.UserID,
		ViewCount: m.ViewCount,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
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
		"_id":        blogObjectID,
		"title":      blog.Title,
		"content":    blog.Content,
		"tags":       blog.Tags,
		"view_count": blog.ViewCount,
		"user_id":    blog.UserID,
		"createdAt":  blog.CreatedAt,
		"updatedAt":  blog.UpdatedAt,
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
        if err != nil {
            return nil, fmt.Errorf("repository: invalid date format '%s', expected YYYY-MM-DD", date)
        }
        start := parsedDate
        end := parsedDate.Add(24 * time.Hour)
        filter["createdAt"] = bson.M{
            "$gte": start,
            "$lt":  end,
        }
    }

    findOptions := options.Find()
    if sort == "popular" {
        findOptions.SetSort(bson.D{{Key: "view_count", Value: -1}})
    } else {
        findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
    }

    cursor, err := r.blogCollection.Find(ctx, filter, findOptions)
    if err != nil {
        return nil, fmt.Errorf("repository: failed to fetch blogs: %w", err)
    }
    defer cursor.Close(ctx)

    var results []blogModel
    if err := cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("repository: failed to decode blogs: %w", err)
    }
	blogs := make([]domain.Blog, 0)

    for _, b := range results {
        blogs = append(blogs, toDomainBlog(b))
    }
    return blogs, nil
}

func (r *blogRepository) DeleteBlog(ctx context.Context, blog *domain.Blog) error {
	objID, err := primitive.ObjectIDFromHex(blog.ID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %w", err)
	}
	_, err = r.blogCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func (r *blogRepository) SearchBlogs(ctx context.Context, tag, date, sort, title, userID string) ([]domain.Blog, error) {
    var andFilters []bson.M

    if tag != "" {
        andFilters = append(andFilters, bson.M{"tags": tag})
    }
    if title != "" {
        andFilters = append(andFilters, bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
    }
    if userID != "" {
        andFilters = append(andFilters, bson.M{"user_id": userID})
    }
    if date != "" {
        if parsedDate, err := time.Parse("2006-01-02", date); err == nil {
            start := parsedDate
            end := parsedDate.Add(24 * time.Hour)
            andFilters = append(andFilters, bson.M{"createdAt": bson.M{"$gte": start, "$lt": end}})
        }
    }

    filter := bson.M{}
    if len(andFilters) > 0 {
        filter["$and"] = andFilters
    }

    findOptions := options.Find()
    if sort == "popular" {
        findOptions.SetSort(bson.D{{Key: "view_count", Value: -1}})
    } else {
        findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
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

func (r *blogRepository) GetPaginated(ctx context.Context, page, limit int, sort, authorID string) ([]domain.Blog, int64, error) {
    filter := bson.M{}
    
    if authorID != "" {
        filter["user_id"] = authorID
    }

    // Count total documents
    total, err := r.blogCollection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count blogs: %w", err)
    }

    // Calculate skip
    skip := (page - 1) * limit

    // Set sort options
    findOptions := options.Find()
    findOptions.SetSkip(int64(skip))
    findOptions.SetLimit(int64(limit))
    
    switch sort {
    case "popular", "views":
        findOptions.SetSort(bson.D{{Key: "view_count", Value: -1}})
    case "recent":
        findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
    default:
        findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
    }

    cursor, err := r.blogCollection.Find(ctx, filter, findOptions)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to fetch blogs: %w", err)
    }
    defer cursor.Close(ctx)

    var results []blogModel
    if err := cursor.All(ctx, &results); err != nil {
        return nil, 0, fmt.Errorf("failed to decode blogs: %w", err)
    }

    var blogs []domain.Blog
    for _, b := range results {
        blogs = append(blogs, toDomainBlog(b))
    }
    
    return blogs, total, nil
}

func (r *blogRepository) IncrementViewCount(ctx context.Context, blogID string) error {
    objID, err := primitive.ObjectIDFromHex(blogID)
    if err != nil {
        return fmt.Errorf("invalid blog ID: %w", err)
    }

    update := bson.M{
        "$inc": bson.M{"view_count": 1},
        "$set": bson.M{"updatedAt": time.Now()},
    }

    _, err = r.blogCollection.UpdateByID(ctx, objID, update)
    if err != nil {
        return fmt.Errorf("failed to increment view count: %w", err)
    }

    return nil
}
