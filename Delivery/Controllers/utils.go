package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAuthenticatedUserID(ctx *gin.Context) (string, bool) {
	val, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return "", false
	}
	userID, ok := val.(string)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return "", false
	}
	return userID, true
}
