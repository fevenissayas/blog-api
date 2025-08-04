package domain

import (
	"context"
	"time"
)

type Blog struct {
	ID             string
	Title          string
    Content        string
	UserID 		   string
	Tags           []string
	ViewCount 		int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) (*Blog, error)
	// GetByUser(ctx context.Context, user *User) (*Blog, error)
	// Delete(ctx context.Context, blog *Blog) (error)
	// Update(ctx context.Context, blogid string, blog*Blog) (*Blog, error)

}
type IBlogUsecase interface {
	Create(ctx context.Context, blog *Blog) error
	// Delete(ctx context.Context, blog *Blog) error
	// Update(ctx context.Context, blogid string, blog*Blog )error
	// Search(ctx context.Context, blogid string) error
	// Filtration(ctx context.Context) error
	// PopulatityTracking(ctx context.Context) error
}
