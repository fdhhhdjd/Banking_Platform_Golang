package gapi

import (
	"context"

	handle "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/handler"
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/auth"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/val"
	"github.com/fdhhhdjd/Banking_Platform_Golang/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *SimpleBankServer) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateCreateUserRequest(req)
	if violations != nil {
		return nil, InvalidArgumentError(violations)
	}

	hashedPassword, err := auth.EncodePassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := database.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	store := db.GetStore()

	user, err := store.CreateUser(c, arg)

	if err != nil {
		errCode := handle.ErrorCode(err)
		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: ConvertUser(user),
	}
	return rsp, nil
}

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, FieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, FieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, FieldViolation("email", err))
	}

	return violations
}
