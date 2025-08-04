package controllers

import (
	domain "blog-api/Domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createBlogRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type updateBlogRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type BlogController struct {
	blogUsecase domain.IBlogUsecase
}

func NewBlogController(blogUsecase domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase: blogUsecase}
}

func (bc *BlogController) CreateBlogHandler(ctx *gin.Context) {
	var req createBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		return
	}

	now := time.Now()
	blog := &domain.Blog{
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		ViewCount: 0,
	}

	if err := bc.blogUsecase.Create(ctx, blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully"})
}

func (bc *BlogController) UpdateBlogHandler(ctx *gin.Context) {

	blogID := ctx.Param("id")

	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		return
	}
	var req updateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	input := domain.UpdateBlogInput{
		BlogID:  blogID,
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}

	updatedBlog, err := bc.blogUsecase.Update(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedBlog)
}
