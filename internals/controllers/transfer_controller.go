package controllers

import (
	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func CreateTransfer(c *gin.Context) error {
	resultTransfer := services.CreateTransfer(c)
	if resultTransfer == nil {
		return nil
	}
	success_response.Created(c, "Create Transfer", resultTransfer)

	return nil
}
