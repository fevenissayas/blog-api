package domain

import "context"

type Comment struct{
	ID 			string
	BlogId 		string
	UserId 		string
	Content 	string
	CreatedAt 	string
}


type ICommentUsecase interface {
    CreateComment(ctx context.Context, blogID, userID, content string) (*Comment, error)
    GetCommentsByBlogID(ctx context.Context, blogID string) ([]Comment, error)
    DeleteComment(ctx context.Context, commentID, userID string) error
}

type ICommentRepository interface {
    Create(ctx context.Context, comment *Comment) (*Comment, error)
    GetByBlogID(ctx context.Context, blogID string) ([]Comment, error)
    Delete(ctx context.Context, commentID, userID string) error
}

type CreateCommentRequest struct {
    Content string `json:"content"`
}