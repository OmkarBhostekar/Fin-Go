package gapi

import (
	"context"
	"database/sql"
	"strings"
	"time"

	db "example.com/simplebank/db/sqlc"
	"example.com/simplebank/pb"
	"example.com/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized: %s", err)
	}

	if authPayload.Username != strings.ToLower(req.GetUsername()) {
		return nil, status.Errorf(codes.PermissionDenied, "can't update other user's profile")
	}

	arg := db.UpdateUserParams{
		Username: strings.ToLower(req.GetUsername()),
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
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
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

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found ")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
