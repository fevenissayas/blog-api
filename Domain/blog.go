package domain

import (
	"context"
	"time"
)

type Blog struct {
	ID        string
	Title     string
	Content   string
	UserID    string
	Tags      []string
	ViewCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) (*Blog, error)
	FindByID(ctx context.Context, BlogID string) (*Blog, error)
	// GetByUser(ctx context.Context, user *User) (*Blog, error)
	DeleteBlog(ctx context.Context, blog *Blog) error
	Update(ctx context.Context, blog *Blog) (*Blog, error)
	Filter(ctx context.Context, tag string, date string, sort string) ([]Blog, error)
}
type IBlogUsecase interface {
	Create(ctx context.Context, blog *Blog) error
	Update(ctx context.Context, input UpdateBlogInput) (*Blog, error)
	DeleteBlog(ctx context.Context, blogID, userID, userRole string) error
	FilterBlogs(ctx context.Context, tag string, date string, sort string) ([]Blog, error)
	// Search(ctx context.Context, blogid string) error
	// Filtration(ctx context.Context) error
	// PopulatityTracking(ctx context.Context) error
}
