package server

import (
	"net/http"
	"os"
	"strconv"
	"time"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	config "github.com/fdhhhdjd/Banking_Platform_Golang/configs"
	routes "github.com/fdhhhdjd/Banking_Platform_Golang/internals/routers"
	util "github.com/fdhhhdjd/Banking_Platform_Golang/internals/utils"
	logger_pkg "github.com/fdhhhdjd/Banking_Platform_Golang/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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

	routes.SetupRoutes(server)

	// 404 - Not Found
	server.NoRoute(NotFound())

	// 500 - Internal Server Error
	server.Use(ServerError)

	server.Run(":" + port)
}

func Pong(c *gin.Context) {
	currentTime := time.Now().Unix()
	signal := "Nguyen Tien Tai"
	message := util.P("Pong %s", signal)

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

func ServerError(c *gin.Context) {
	c.Next()

	log := logrus.New()

	if len(c.Errors) > 0 {
		err := c.Errors[0].Err
		status := http.StatusInternalServerError

		logEntry := log.WithFields(logrus.Fields{
			"status":  status,
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"time":    time.Now(),
			"request": c.Request.RequestURI,
		})
		logEntry.Error(err)
		errorResponse := error_response.InternalServerError("Internal Server Error")
		c.JSON(errorResponse.Status, errorResponse)
	}
}
