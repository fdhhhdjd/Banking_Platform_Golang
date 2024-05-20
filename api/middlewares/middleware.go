package middlewares

import (
	"strings"
	"time"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/auth"
	"github.com/gin-gonic/gin"
)

type Payload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			errorResponse := error_response.UnauthorizedError("")
			c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
			errorResponse := error_response.UnauthorizedError("")
			c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
			return
		}

		accessToken := fields[1]

		JwtMaker, err := auth.GetJWTMaker()
		if err != nil {
			errorResponse := error_response.UnauthorizedError("")
			c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
			return
		}

		payload, err := JwtMaker.VerifyToken(accessToken)

		c.Set("info_user", payload)

		c.Next()
	}
}
