package usecases

import (
	domain "blog-api/Domain"
	"context"
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
    _,err := bu.blogRepository.Create(ctx,blog)
	return err
}
