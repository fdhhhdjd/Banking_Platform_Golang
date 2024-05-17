package services

import (
	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
	"github.com/gin-gonic/gin"
)

func GetAllAccount(c *gin.Context) []models.Account {
	type listAccountRequest struct {
		PageID   int32 `form:"page_id" binding:"required,min=1"`
		PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	}
	var req listAccountRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	ctx := c.Request.Context()
	arg := database.ListAccountsParams{
		Owner:  "pitdjj", // Example user dummy
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	store := db.GetStore()

	accounts, err := store.ListAccounts(ctx, arg)
	if err != nil {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	if len(accounts) == 0 {
		return []models.Account{}
	}

	var result []models.Account
	for _, account := range accounts {
		result = append(result, models.Account(account))
	}

	return result
}
