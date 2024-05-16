package controllers

import (
	"net/http"

	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	user_services "github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) error {
	getUsers := user_services.GetAllUsers(c)
	if getUsers == nil {
		return nil
	}

	successResponse := success_response.NewSuccessResponse("OK Good", http.StatusOK, http.StatusText(http.StatusOK), nil, getUsers)
	successResponse.Send(c)
	return nil
}
