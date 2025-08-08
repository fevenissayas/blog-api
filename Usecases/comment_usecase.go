package usecases

import (
    "blog-api/Domain"
    "context"
    "time"
)

type CommentUsecase struct {
    CommentRepository domain.ICommentRepository
}

func NewCommentUsecase(repo domain.ICommentRepository) domain.ICommentUsecase {
    return &CommentUsecase{CommentRepository: repo}
}

func (u *CommentUsecase) CreateComment(ctx context.Context, blogID, userID, content string) (*domain.Comment, error) {
    comment := &domain.Comment{
        ID:        "",
        BlogId:    blogID,
        UserId:    userID,
        Content:   content,
        CreatedAt: time.Now().Format(time.RFC3339),
    }
    return u.CommentRepository.Create(ctx, comment)
}

func (u *CommentUsecase) GetCommentsByBlogID(ctx context.Context, blogID string) ([]domain.Comment, error) {
    return u.CommentRepository.GetByBlogID(ctx, blogID)
}

func (u *CommentUsecase) DeleteComment(ctx context.Context, commentID, userID string) error {
    return u.CommentRepository.Delete(ctx, commentID, userID)
}