package controllers

import (
	domain "blog-api/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase domain.IAuthUsecase
}

func NewAuthController(authUsecase domain.IAuthUsecase) *AuthController {
	return &AuthController{authUsecase: authUsecase}
}

func (a *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	// First try to get refresh token from cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		// Fallback to Authorization header for backward compatibility
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
			return
		}

		// Parse Bearer token format
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			refreshToken = authHeader[7:]
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
	}

	tokens, err := a.authUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		// Clear invalid refresh token cookie
		ctx.SetCookie(
			"refresh_token", // name
			"",              // value (empty to clear)
			-1,              // maxAge (negative to expire immediately)
			"/",             // path
			"",              // domain
			false,           // secure
			true,            // httpOnly
		)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set new refresh token as HTTP-only cookie
	ctx.SetCookie(
		"refresh_token",           // name
		tokens.RefreshToken,       // value
		7*24*60*60,               // maxAge (7 days in seconds)
		"/",                      // path
		"",                       // domain (empty means current domain)
		false,                    // secure (set to true in production with HTTPS)
		true,                     // httpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"message":      "Token refreshed successfully",
	})
}
