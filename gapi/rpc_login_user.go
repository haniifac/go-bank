package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/pb"
	"github.com/haniifac/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, status.Errorf(codes.NotFound, "user not found")
		default:
			return nil, status.Errorf(codes.Internal, "cannot login: %s", err)
		}
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect password")
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	args := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		IpAddress:    "",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(server.config.RefreshTokenDuration),
		IsBlocked:    false,
	}

	_, err = server.store.CreateSession(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user session: %s", err)
	}

	res := &pb.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			Username:  user.Username,
			FullName:  user.Fullname,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.PasswordUpdatedAt),
		},
	}

	return res, nil
}
