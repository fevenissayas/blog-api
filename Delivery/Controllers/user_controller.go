package controllers

import (
	domain "blog-api/Domain"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(userUsecase domain.IUserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

type RegisterRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Bio            string `json:"bio"`             //optional
	ProfilePicture string `json:"profile_picture"` //optional
	ContactInfo    string `json:"contact_info"`    //optional
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) RegisterHandler(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	if !govalidator.IsEmail(req.Email) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid email format"})
		return
	}
	//convert DTO to domain.User
	user := &domain.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       req.Password,
		Bio:            req.Bio,
		ProfilePicture: req.ProfilePicture,
		ContactInfo:    req.ContactInfo,
		Role:           domain.RoleUser,
		IsVerified:     false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := c.userUsecase.Register(ctx, user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (c *UserController) LoginHandler(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}
	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	tokens, err := c.userUsecase.Login(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (c *UserController) LogoutHandler(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		return
	}

	if err := c.userUsecase.Logout(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
