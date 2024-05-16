package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	config "github.com/fdhhhdjd/Banking_Platform_Golang/configs"
	routes "github.com/fdhhhdjd/Banking_Platform_Golang/internals/routers"
	logger_pkg "github.com/fdhhhdjd/Banking_Platform_Golang/pkg"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	logger_pkg.InitLog()
	config.LoadConfig()
}

func Server() {
	server := gin.New()
	// Manually add middleware
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	port := os.Getenv("PORT")

	if port == "" {
		port = strconv.Itoa(config.AppConfig.Server.Port)
	}
	server.GET("/ping", Pong)

	// 404
	server.NoRoute(NotFound())

	// 500 - Internal Server Error
	server.Use(ServerError())

	routes.SetupRoutes(server)

	server.Run(":" + port)
}

var ErrDivideByZero = errors.New("divide by zero")

func Pong(c *gin.Context) {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	message := fmt.Sprintf("Pong %s Hello Tai Dev", currentTime)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": message,
		"time":    currentTime,
	})
}
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
	}
}

func ServerError() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorResponse := error_response.InternalServerError("Internal Server Error")
		c.JSON(errorResponse.Status, errorResponse)
	}
}
