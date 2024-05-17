package controllers

import (
	"net/http"

	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func GetAllAccount(c *gin.Context) error {
	getAccount := services.GetAllAccount(c)
	if getAccount == nil {
		return nil
	}

	successResponse := success_response.NewSuccessResponse("OK Good", http.StatusOK, http.StatusText(http.StatusOK), nil, getAccount)
	successResponse.Send(c)
	return nil
}
