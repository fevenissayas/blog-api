package usecases

import (
	domain "blog-api/Domain"
	"context"
	"fmt"
	"time"
)

type BlogUsecase struct {
	blogRepository domain.IBlogRepository
	aiService domain.AiService
}

func NewBlogUseCase(blogRepo domain.IBlogRepository, aiservice domain.AiService) domain.IBlogUsecase {
	return &BlogUsecase{
		blogRepository: blogRepo,
		aiService: aiservice,
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

func (bu *BlogUsecase) SearchBlogs(ctx context.Context, tag, date, sort, title, userID string) ([]domain.Blog, error) {
    return bu.blogRepository.SearchBlogs(ctx, tag, date, sort, title, userID)
}
func (bu *BlogUsecase) GetSuggestion(req domain.AiSuggestionRequest)(string , error){
	return bu.aiService.Getsuggestion(req)
}

func (bu *BlogUsecase) GetBlogs(ctx context.Context, page, limit int, sort, authorID string) (*domain.PaginatedBlogsResponse, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    blogs, total, err := bu.blogRepository.GetPaginated(ctx, page, limit, sort, authorID)
    if err != nil {
        return nil, fmt.Errorf("failed to get blogs: %w", err)
    }

    totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division
    
    response := &domain.PaginatedBlogsResponse{
        Blogs:      blogs,
        Page:       page,
        Limit:      limit,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }

    return response, nil
}

func (bu *BlogUsecase) GetByIDAndIncrementViews(ctx context.Context, blogID string) (*domain.Blog, error) {
    // First get the blog
    blog, err := bu.blogRepository.FindByID(ctx, blogID)
    if err != nil {
        return nil, fmt.Errorf("blog not found: %w", err)
    }

    // Increment view count
    err = bu.blogRepository.IncrementViewCount(ctx, blogID)
    if err != nil {
        // Log error but don't fail the request
        fmt.Printf("Warning: failed to increment view count for blog %s: %v\n", blogID, err)
    } else {
        // Update the blog object with incremented count
        blog.ViewCount++
    }

    return blog, nil
}
