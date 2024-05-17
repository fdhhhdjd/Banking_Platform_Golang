package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	//* All
	router.GET("/users", handle.AsyncHandler(controllers.GetUsers))

	//* Create
	router.POST("/users/register", handle.AsyncHandler(controllers.Register))

}
