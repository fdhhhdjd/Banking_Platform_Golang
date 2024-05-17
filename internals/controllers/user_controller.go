package controllers

import (
	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	user_services "github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) error {
	getUsers := user_services.GetAllUsers(c)
	if getUsers == nil {
		return nil
	}
	success_response.Ok(c, "Get All", getUsers)
	return nil
}

func Register(c *gin.Context) error {
	registerUser := user_services.RegisterUser(c)
	if registerUser == nil {
		return nil
	}
	success_response.Created(c, "Register", registerUser)
	return nil
}

func Login(c *gin.Context) error {
	registerUser := user_services.RegisterUser(c)
	if registerUser == nil {
		return nil
	}
	success_response.Created(c, "Register", registerUser)
	return nil
}
