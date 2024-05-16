package handle

import (
	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	"github.com/gin-gonic/gin"
)

func AsyncHandler(fn func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			c.Error(err)
			errorResponse := error_response.InternalServerError(err.Error())
			c.JSON(errorResponse.Status, errorResponse)
		}
	}
}
