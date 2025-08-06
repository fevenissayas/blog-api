package usecases

import (
	domain "blog-api/Domain"
	"context"
	"fmt"
	"time"
)

type BlogUsecase struct {
	blogRepository domain.IBlogRepository
}

func NewBlogUseCase(blogRepo domain.IBlogRepository) domain.IBlogUsecase {
	return &BlogUsecase{
		blogRepository: blogRepo,
	}
}

func (bu *BlogUsecase) Create(ctx context.Context, blog *domain.Blog) error {

	if blog.Title == "" || blog.Content == "" || len(blog.Tags) == 0 {
		return fmt.Errorf("invalid input: title/content/tags must not be empty")
	}
	_, err := bu.blogRepository.Create(ctx, blog)

	if err != nil {
		return fmt.Errorf("failed to create blog: %w", err)
	}
	return nil
}

func (bu *BlogUsecase) Update(ctx context.Context, input domain.UpdateBlogInput) (*domain.Blog, error) {
	blog, err := bu.blogRepository.FindByID(ctx, input.BlogID)
	if err != nil {
		return nil, fmt.Errorf("blog not found")
	}

	if blog.UserID != input.UserID {
		return nil, fmt.Errorf("unauthorized: you are not the author")
	}

	if input.Title != "" {
		blog.Title = input.Title
	}
	if input.Content != "" {
		blog.Content = input.Content
	}
	if len(input.Tags) != 0 {
		blog.Tags = input.Tags
	}

	blog.UpdatedAt = time.Now()

	updatedBlog, err := bu.blogRepository.Update(ctx, blog)
	if err != nil {
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}

	return updatedBlog, nil
}

func (bu *BlogUsecase) FilterBlogs(ctx context.Context, tag string, date string, sort string) ([]domain.Blog, error) {
	return bu.blogRepository.Filter(ctx, tag, date, sort)
}

func (bu *BlogUsecase) DeleteBlog(ctx context.Context, blogID, userID, userRole string) error {
	blog, err := bu.blogRepository.FindByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found")
	}
	if blog.UserID != userID && userRole != "admin" {
		return fmt.Errorf("unauthorized: only the author or admin can delete this blog")
	}
	return bu.blogRepository.DeleteBlog(ctx, blog)
}

