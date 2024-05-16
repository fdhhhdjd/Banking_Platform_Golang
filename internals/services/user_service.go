package services

import (
	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) []models.User {
	users := []models.User{
		{ID: 1, Name: "Nguyen Tien Tai", Email: "tai@example.com"},
	}

	if len(users) > 0 {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	return users
}
