package controllers

import (
	"blog-api/Domain"
	"blog-api/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
type UserController struct{
     UserUsecase *usecases.UserUsecase    
} 
func NewUserController(uuc *usecases.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: uuc,
	}
}
func (uc *UserController) Register (c *gin.Context){
	var user domain.User
    if err := c.BindJSON(&user); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"invalid input"})
		return
	}
	// User Validation
	validate := validator.New()
    if err := validate.Struct(user); err != nil {
        c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
        return
    }

	if err := uc.UserUsecase.Register(&user); err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})	
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message":"Succesfully Registered User"})
}

func (uc *UserController) Login(c *gin.Context){
	var req domain.LoginRequest
	if err := c.BindJSON(&req); err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := uc.UserUsecase.AuthenticateUser(req.UsernameOrEmail, req.Password)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": token})

}