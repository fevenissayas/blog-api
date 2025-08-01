package router

import (
	controller "blog-api/Delivery/Controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(c *controller.UserController) *gin.Engine{
	router := gin.Default()

	router.POST("/register",c.Register)
    router.POST("/login", c.Login)
	return router
}
