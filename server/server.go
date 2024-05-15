package server

import (
	"net/http"
	"os"
	"strconv"

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
	router := gin.New()

	port := os.Getenv("PORT")

	if port == "" {
		port = strconv.Itoa(config.AppConfig.Server.Port)
	}
	router.GET("/", homePage)

	routes.SetupRoutes(router)

	router.Run(":" + port)
}

func homePage(c *gin.Context) {
	// c.String(http.StatusOK, "This is my home page")
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Hello Tai Dev",
	})
}
