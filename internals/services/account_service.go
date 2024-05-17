package services

import (
	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
	"github.com/gin-gonic/gin"
)

func GetAllAccount(c *gin.Context) []models.Account {

	var req models.ListAccountRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	arg := database.ListAccountsParams{
		Owner:  "Taidev", // Example user dummy
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	store := db.GetStore()

	accounts, err := store.ListAccounts(c, arg)
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
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

func CreateAccount(c *gin.Context) *models.Account {
	var req models.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	store := db.GetStore()

	arg := database.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	accounts, err := store.CreateAccount(c, arg)
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	newAccount := models.Account{
		ID:        accounts.ID,
		Owner:     accounts.Owner,
		Balance:   accounts.Balance,
		Currency:  accounts.Currency,
		CreatedAt: accounts.CreatedAt,
	}

	return &newAccount
}

func GetAccount(c *gin.Context) *models.Account {
	var req models.GetAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}
	store := db.GetStore()

	account, err := store.GetAccount(c, req.ID)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	newAccount := models.Account{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}

	return &newAccount
}
