package services

import (
	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
	util "github.com/fdhhhdjd/Banking_Platform_Golang/internals/utils"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) []models.User {
	users := []models.User{
		{ID: 1, Name: "Nguyen Tien Tai", Email: "tai@example.com"},
	}

	if len(users) == 0 {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	return users
}

func RegisterUser(c *gin.Context) *database.User {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	hashedPassword, err := util.EncodePassword(req.Password)

	if err != nil {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	arg := database.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	store := db.GetStore()

	user, err := store.CreateUser(c, arg)

	if err != nil {
		errCode := handle.ErrorCode(err)
		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
			errorResponse := error_response.ForbiddenError(errCode)
			c.JSON(errorResponse.Status, errorResponse)
			return nil
		}
		errorResponse := error_response.InternalServerError("")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}
	return &user
}

func LoginUser(c *gin.Context) *models.LoginUserResponse {
	var req models.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := error_response.NotFoundError("Not Found")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	// store := db.GetStore()

	// user, err := store.GetUser(c, req.Username)

	// if err != nil {
	// 	if err != nil {
	// 		errCode := handle.ErrorCode(err)
	// 		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
	// 			errorResponse := error_response.ForbiddenError(errCode)
	// 			c.JSON(errorResponse.Status, errorResponse)
	// 			return nil
	// 		}
	// 		errorResponse := error_response.InternalServerError("")
	// 		c.JSON(errorResponse.Status, errorResponse)
	// 		return nil
	// 	}
	// }

	return nil
}
