package controllers

import (
	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func GetAllAccount(c *gin.Context) error {
	getAccount := services.GetAllAccount(c)
	if getAccount == nil {
		return nil
	}
	success_response.Ok(c, "Get all account", getAccount)

	return nil
}

func CreateAccount(c *gin.Context) error {
	createAccount := services.CreateAccount(c)
	if createAccount == nil {
		return nil
	}
	success_response.Created(c, "Create Account", createAccount)

	return nil
}
