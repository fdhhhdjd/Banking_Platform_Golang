package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	controller "github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/gin-gonic/gin"
)

func AccountRoutes(router *gin.Engine) {
	//* Get all
	router.GET("/account/all", handle.AsyncHandler(controller.GetAllAccount))

	//* Created account
	router.POST("/account", handle.AsyncHandler(controller.CreateAccount))
}
