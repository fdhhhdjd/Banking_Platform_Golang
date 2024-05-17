package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/gin-gonic/gin"
)

func EntriesRoutes(router *gin.Engine) {
	//* Entries
	router.GET("/entries", handle.AsyncHandler(controllers.GetAllEntries))
}
