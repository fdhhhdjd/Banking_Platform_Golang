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
	//* Get Cookie
	resultRefetch, err := c.Cookie(constants.KeyRefetchToken)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Get Info for accessToken
	payload, exists := c.Get("info_user")

	if !exists {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}
	payloadPtr, ok := payload.(*auth.Payload)

	if !ok {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	users := []models.User{
		{ID: payloadPtr.ID, Name: payloadPtr.Username, Email: "tai@gmail.com", RefetchToken: resultRefetch},
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

	//* Check Password
	DecodePassword := auth.DecodePassword(req.Password, user.HashedPassword)
	if DecodePassword != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Handle Secret key token
	JwtMaker, err := auth.GetJWTMaker()
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Handle create token
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

	accessToken, accessPayload, err := JwtMaker.CreateToken(
		user.Username,
		user.Role,
		accessTokenDuration,
	)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	refreshToken, refreshPayload, err := JwtMaker.CreateToken(
		user.Username,
		user.Role,
		refetchTokenDuration,
	)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Save session account login
	session, err := store.CreateSession(c, database.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		errorResponse := error_response.InternalServerError("Internal Server Error")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Set cookie with access token
	nodeEnv := os.Getenv("ENV")
	domain := constants.HOST
	secure := nodeEnv != constants.DEV

	if nodeEnv != constants.DEV {
		hostWithPort := c.Request.Host
		parts := strings.Split(hostWithPort, ":")
		domain = parts[0]
	}

	//* Set cookie with access token
	c.SetCookie(constants.KeyRefetchToken, refreshToken, int(refetchTokenDuration.Seconds()), "/", domain, secure, true)

	rsp := models.LoginUserResponse{
		SessionID:             session.ID,
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

func RenewToken(c *gin.Context) *models.RenewAccessTokenResponse {
	//* Check Input
	var req models.RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Handle check RefetchToken
	JwtMaker, err := auth.GetJWTMaker()
	if err != nil {
		errorResponse := error_response.UnauthorizedError("")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return nil
	}

	RefreshPayload, err := JwtMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Get info section Db
	store := db.GetStore()
	session, err := store.GetSession(c, RefreshPayload.ID)

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

	//* Check account have been block yet
	if session.IsBlocked {
		errorResponse := error_response.UnauthorizedError("")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Check account have in the same with username session
	if session.Username != RefreshPayload.Username {
		errorResponse := error_response.UnauthorizedError("")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Check refetch token not yet
	if session.RefreshToken != req.RefreshToken {
		errorResponse := error_response.UnauthorizedError("")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Check session have been expired yet
	if time.Now().After(session.ExpiresAt) {
		errorResponse := error_response.UnauthorizedError("")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return nil
	}

	//* Create accessToken new
	accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}
	accessToken, accessPayload, err := JwtMaker.CreateToken(
		RefreshPayload.Username,
		RefreshPayload.Role,
		accessTokenDuration,
	)

	if err != nil {
		errorResponse := error_response.BadRequestError("Bad Request")
		c.JSON(errorResponse.Status, errorResponse)
		return nil
	}

	rsp := models.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return &rsp
}
