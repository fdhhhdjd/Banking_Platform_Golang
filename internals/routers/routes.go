package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//* Users
	UserRoutes(router)

	//* Entries
	EntriesRoutes(router)

	//* Account
	AccountRoutes(router)
}
