package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	"github.com/fdhhhdjd/Banking_Platform_Golang/api/middlewares"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	//* Register
	router.POST("/users/register", handle.AsyncHandler(controllers.Register))

	//* Login
	router.POST("/users/login", handle.AsyncHandler(controllers.Login))

	//* All
	protected := router.Group("/", middlewares.AuthMiddleware())
	{
		protected.GET("/users", handle.AsyncHandler(controllers.GetUsers))
	}
}
