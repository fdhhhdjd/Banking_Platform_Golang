package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	controller "github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//* Users
	router.GET("/users", handle.AsyncHandler(controller.GetUsers))

	//* Entries
	router.GET("/entries", handle.AsyncHandler(controller.GetAllEntries))

	//* Account
	router.GET("/account", handle.AsyncHandler(controller.GetAllAccount))

}
