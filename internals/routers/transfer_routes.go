package routes

import (
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	controller "github.com/fdhhhdjd/Banking_Platform_Golang/internals/controllers"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/helpers"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func TransferRoutes(router *gin.Engine) {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", helpers.ValidCurrency)
	}

	//* Created
	router.POST("/transfer", handle.AsyncHandler(controller.CreateTransfer))
}
