package tests

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func sampleErrorFunction(c *gin.Context) error {
	return errors.New("Something went wrong")
}

func handlerWithAsync(c *gin.Context) error {
	return sampleErrorFunction(c)
}
