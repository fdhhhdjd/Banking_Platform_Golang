package gapi

import (
	"context"
	"database/sql"
	"time"

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

func (server *SimpleBankServer) UpdateUser(c context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := ValidateUpdateUserRequest(req)
	if violations != nil {
		return nil, InvalidArgumentError(violations)
	}

	arg := database.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := auth.EncodePassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	store := db.GetStore()

	user, err := store.UpdateUser(c, arg)

	if err != nil {
		errCode := handle.ErrorCode(err)
		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: ConvertUser(user),
	}
	return rsp, nil
}

func ValidateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, FieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, FieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, FieldViolation("email", err))
		}
	}

	return violations
}
