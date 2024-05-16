package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	controller "github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/gin-gonic/gin"
)

// Hàm mẫu mô phỏng trả về lỗi

func SetupRoutes(router *gin.Engine) {
	//* Users
	router.GET("/users", handle.AsyncHandler(controller.GetUsers))

	//* Entries
	router.GET("/entries", handle.AsyncHandler(controller.GetAllEntries))

}
