package handle

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func AsyncHandler(fn func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			c.Error(err)
			c.Next()
		}
	}
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	}

	var pgErr *pgconn.PgError
	log.Println(pgErr)
	if errors.As(err, &pgErr) && pgErr != nil {
		return pgErr.Code
	}

	return ""
}
