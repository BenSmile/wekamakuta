package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/db/util"
	"github.com/bensmile/wekamakuta/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get the user : %s", err)

	}

	if err := util.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "invalid password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the token : %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the refresh token : %s", err)
	}

	sessionParam := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, sessionParam)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}
