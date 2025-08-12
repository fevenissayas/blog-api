package controllers

import (
	domain "blog-api/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
    commentUsecase domain.ICommentUsecase
}

func NewCommentController(commentUsecase domain.ICommentUsecase) *CommentController {
    return &CommentController{commentUsecase: commentUsecase}
}

// Create a new comment for a blog
func (cc *CommentController) CreateComment(ctx *gin.Context) {
    blogID := ctx.Param("id")
    var req domain.CreateCommentRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    userID, ok := getAuthenticatedUserID(ctx)
    if !ok {
        return
    }
    comment, err := cc.commentUsecase.CreateComment(ctx.Request.Context(), blogID, userID, req.Content)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, comment)
}

// Get all comments for a blog
func (cc *CommentController) GetComments(ctx *gin.Context) {
    blogID := ctx.Param("id")
    comments, err := cc.commentUsecase.GetCommentsByBlogID(ctx.Request.Context(), blogID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, comments)
}

// Delete a comment
func (cc *CommentController) DeleteComment(ctx *gin.Context) {
    commentID := ctx.Param("commentID")
    userID, ok := getAuthenticatedUserID(ctx)
    if !ok {
        return
    }
    err := cc.commentUsecase.DeleteComment(ctx.Request.Context(), commentID, userID)
    if err != nil {
        ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}