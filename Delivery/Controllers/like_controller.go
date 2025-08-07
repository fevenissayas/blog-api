package controllers

import (
	domain "blog-api/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	likeUsecase domain.IBlogLikeUsecase
}

func NewLikeController(likeUsecase domain.IBlogLikeUsecase) *LikeController {
	return &LikeController{
		likeUsecase: likeUsecase,
	}
}

func (c *LikeController) LikeBlogHandler(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "blog id is required"})
		return
	}

	err := c.likeUsecase.LikeBlog(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "blog liked"})
}

func (c *LikeController) UnlikeBlogHandler(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "blog id is required"})
		return
	}

	err := c.likeUsecase.UnlikeBlog(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "blog unliked"})
}

func (c *LikeController) GetLikeCountHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "blog id is required"})
		return
	}

	count, err := c.likeUsecase.GetLikeCount(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"like_count": count})
}

func (c *LikeController) IsBlogLikedHandler(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "blog id is required"})
		return
	}

	isLiked, err := c.likeUsecase.IsBlogLiked(id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"is_liked": isLiked})
}
