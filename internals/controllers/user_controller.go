package controllers

import (
	"net/http"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	user_services "github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := user_services.GetAllUsers()
	if err != nil {
		errorResponse := error_response.InternalServerError(err.Error())
		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	successResponse := success_response.NewSuccessResponse("OK Good", http.StatusOK, http.StatusText(http.StatusOK), nil, users)
	successResponse.Send(c)
}
