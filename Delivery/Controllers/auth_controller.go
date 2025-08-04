package controllers

import (
	domain "blog-api/Domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase domain.IAuthUsecase
}

func NewAuthController(authUsecase domain.IAuthUsecase) *AuthController {
	return &AuthController{authUsecase: authUsecase}
}

func (a *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	refreshToken := parts[1]

	tokens, err := a.authUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}
