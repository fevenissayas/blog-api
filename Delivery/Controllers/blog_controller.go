package controllers

import (
	domain "blog-api/Domain"
	"net/http"
	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUsecase domain.IBlogUsecase
}

func NewBlogController(blogUsecase domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase: blogUsecase}
}
func (bc *BlogController) Create (ctx *gin.Context){
    var blog domain.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
    if err := bc.blogUsecase.Create(ctx, &blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully"})
}
