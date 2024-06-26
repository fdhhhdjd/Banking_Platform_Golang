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
	loginUser := user_services.LoginUser(c)
	if loginUser == nil {
		return nil
	}
	success_response.Ok(c, "Login", loginUser)
	return nil
}

func RenewToken(c *gin.Context) error {
	RenewToken := user_services.RenewToken(c)
	if RenewToken == nil {
		return nil
	}
	success_response.Ok(c, "RenewToken", RenewToken)
	return nil
}
