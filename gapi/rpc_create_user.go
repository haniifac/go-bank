package gapi

import (
	"context"

	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/pb"
	"github.com/haniifac/simplebank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "cannot hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
		Username:       req.GetUsername(),
		Fullname:       req.GetFullName(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %s", pqErr.Message)
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %s", err)

	}

	pbUser := &pb.User{
		Username:  user.Username,
		FullName:  user.Fullname,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.PasswordUpdatedAt),
	}

	res := &pb.CreateUserResponse{
		User: pbUser,
	}

	return res, nil
}
