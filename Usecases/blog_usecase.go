package usecases

import (
	domain "blog-api/Domain"
	"context"
	"fmt"
)

type BlogUsecase struct {
	blogRepository   domain.IBlogRepository
}

func NewBlogUseCase(blogRepo domain.IBlogRepository) domain.IBlogUsecase {
	return &BlogUsecase{
		blogRepository:   blogRepo,
	}
}

func (bu *BlogUsecase) Create (ctx context.Context,blog *domain.Blog) error{

		if blog.Title == "" || blog.Content == "" || len(blog.Tags)== 0 {
		return domain.ErrInvalidInput
	}
    _,err := bu.blogRepository.Create(ctx,blog)

		if err != nil {
		return fmt.Errorf("failed to create blog: %w", err)
	}
	return nil
}
