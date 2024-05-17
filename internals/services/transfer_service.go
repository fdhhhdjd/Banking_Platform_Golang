package services

import (
	"database/sql"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
	"github.com/gin-gonic/gin"
)

func CreateTransfer(c *gin.Context) *database.TransferTxResult {
	var req models.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Check validate account from
	if !ValidateAccount(c, req.FromAccountID, req.Currency) {
		return nil
	}

	//* Check validate account to
	if !ValidateAccount(c, req.ToAccountID, req.Currency) {
		return nil
	}

	store := db.GetStore()

	arg := database.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := store.TransferTx(c, arg)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")

		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	return &result
}

func ValidateAccount(c *gin.Context, accountID int64, currency string) bool {
	store := db.GetStore()
	account, err := store.GetAccount(c, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := error_response.InternalServerError("")
			c.JSON(errorResponse.Status, errorResponse)
			return false
		}

		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return false
	}

	if account.Currency != currency {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return false
	}

	return true
}
