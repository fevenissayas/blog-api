package controllers

import (
	domain "blog-api/Domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUsecase domain.IBlogUsecase
}

func NewBlogController(blogUsecase domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase: blogUsecase}
}
func (bc *BlogController) Create(ctx *gin.Context) {
	var blog domain.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(string)
	if !ok || userID == ""{
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	blog.UserID = userID

	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now
	blog.ViewCount = 0

	if err := bc.blogUsecase.Create(ctx, &blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully"})
}
