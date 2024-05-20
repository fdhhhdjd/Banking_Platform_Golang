package services

import (
	"os"
	"strings"
	"time"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/auth"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) []models.User {
	users := []models.User{
		{ID: 1, Name: "Nguyen Tien Tai", Email: "tai@example.com"},
	}

	if len(users) == 0 {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	return users
}

func RegisterUser(c *gin.Context) *database.User {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	hashedPassword, err := auth.EncodePassword(req.Password)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
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
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	store := db.GetStore()

	user, err := store.GetUser(c, req.Username)

	if err != nil {
		errCode := handle.ErrorCode(err)
		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
			errorResponse := error_response.ForbiddenError(errCode)
			c.JSON(errorResponse.Status, errorResponse)
			return nil
		}
		errorResponse := error_response.InternalServerError("Internal Server Error")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	JwtMaker, err := auth.GetJWTMaker()
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	refetchTokenDurationStr := os.Getenv("REFRESH_TOKEN_DURATION")

	accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	refetchTokenDuration, err := time.ParseDuration(refetchTokenDurationStr)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	accessToken, accessPayload, err := JwtMaker.CreateToken(
		user.Username,
		user.Role,
		accessTokenDuration,
	)

	if err != nil {
		errorResponse := error_response.InternalServerError("Internal Server Error")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	refreshToken, refreshPayload, err := JwtMaker.CreateToken(
		user.Username,
		user.Role,
		refetchTokenDuration,
	)

	if err != nil {
		errorResponse := error_response.InternalServerError("Internal Server Error")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	nodeEnv := os.Getenv("ENV")
	domain := constants.HOST
	secure := nodeEnv != constants.DEV

	if nodeEnv != constants.DEV {
		hostWithPort := c.Request.Host
		parts := strings.Split(hostWithPort, ":")
		domain = parts[0]
	}

	// Set cookie with access token
	c.SetCookie(constants.KeyRefetchToken, refreshToken, int(refetchTokenDuration.Seconds()), "/", domain, secure, true)

	rsp := models.LoginUserResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: models.UserResponse{
			Username:          user.Username,
			Email:             user.Email,
			FullName:          user.FullName,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
	}

	return &rsp
}
