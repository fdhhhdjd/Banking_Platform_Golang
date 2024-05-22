package gapi

import (
	"context"
	"os"
	"time"

	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/auth"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *SimpleBankServer) LoginUser(c context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	store := db.GetStore()

	user, err := store.GetUser(c, req.Username)

	if err != nil {
		errCode := handle.ErrorCode(err)
		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
			return nil, status.Errorf(codes.Internal, "User: %s", err)
		}
	}

	//* Check Password
	DecodePassword := auth.DecodePassword(req.Password, user.HashedPassword)
	if DecodePassword != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Password: %s", err)

	}

	//* Handle Secret key token
	JwtMaker, err := auth.GetJWTMaker()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "JWT: %s", err)

	}

	//* Handle create token
	accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	refetchTokenDurationStr := os.Getenv("REFRESH_TOKEN_DURATION")

	accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "AccessToken: %s", err)
	}

	refetchTokenDuration, err := time.ParseDuration(refetchTokenDurationStr)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "RefetchToken: %s", err)
	}

	accessToken, accessPayload, err := JwtMaker.CreateToken(
		user.Username,
		user.Role,
		accessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "AccessToken: %s", err)
	}

	refreshToken, refreshPayload, err := JwtMaker.CreateToken(
		user.Username,
		user.Role,
		refetchTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "RefetchToken: %s", err)
	}

	//* Save session account login
	mtdt := server.extractMetadata(c)
	session, err := store.CreateSession(c, database.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Session: %s", err)

	}

	rsp := &pb.LoginUserResponse{
		User:                  ConvertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}
