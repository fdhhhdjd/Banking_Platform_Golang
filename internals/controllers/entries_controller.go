package controllers

import (
	"net/http"

	success_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/success"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/services"
	"github.com/gin-gonic/gin"
)

func GetAllEntries(c *gin.Context) error {
	successResponse := success_response.NewSuccessResponse("OK Good", http.StatusOK, http.StatusText(http.StatusOK), nil, services.GetAllEntries())
	successResponse.Send(c)
	return nil
}
