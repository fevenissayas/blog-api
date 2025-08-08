package repositories

import (
    "blog-api/Domain"
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type commentModel struct {
    ID        primitive.ObjectID `bson:"_id"`
    BlogID    string             `bson:"blog_id"`
    UserID    string             `bson:"user_id"`
    Content   string             `bson:"content"`
    CreatedAt time.Time          `bson:"created_at"`
}

func toDomainComment(m commentModel) domain.Comment {
    return domain.Comment{
        ID:        m.ID.Hex(),
        BlogId:    m.BlogID,
        UserId:    m.UserID,
        Content:   m.Content,
        CreatedAt: m.CreatedAt.Format(time.RFC3339),
    }
}

type CommentRepository struct {
    collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) domain.ICommentRepository {
    return &CommentRepository{collection: db.Collection("comments")}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
    objID := primitive.NewObjectID()
    comment.ID = objID.Hex()
    doc := bson.M{
        "_id":        objID,
        "blog_id":    comment.BlogId,
        "user_id":    comment.UserId,
        "content":    comment.Content,
        "created_at": comment.CreatedAt,
    }
    _, err := r.collection.InsertOne(ctx, doc)
    if err != nil {
        return nil, fmt.Errorf("failed to insert comment: %w", err)
    }
    return comment, nil
}

func (r *CommentRepository) GetByBlogID(ctx context.Context, blogID string) ([]domain.Comment, error) {
    filter := bson.M{"blog_id": blogID}
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch comments: %w", err)
    }
    defer cursor.Close(ctx)

    var results []commentModel
    if err := cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode comments: %w", err)
    }

    var comments []domain.Comment
    for _, c := range results {
        comments = append(comments, toDomainComment(c))
    }
    return comments, nil
}

func (r *CommentRepository) Delete(ctx context.Context, commentID, userID string) error {
    objID, err := primitive.ObjectIDFromHex(commentID)
    if err != nil {
        return fmt.Errorf("invalid comment ID: %w", err)
    }
    filter := bson.M{"_id": objID, "user_id": userID}
    res, err := r.collection.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("failed to delete comment: %w", err)
    }
    if res.DeletedCount == 0 {
        return fmt.Errorf("comment not found or not owned by user")
    }
    return nil
}