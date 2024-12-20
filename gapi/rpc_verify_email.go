package gapi

import (
	"context"

	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/pb"
	"github.com/bensmile/wekamakuta/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {

	violations := ValidateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VeirfyEmailTx(ctx, db.VeirfyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerify: txResult.User.IsEmailVerified,
	}

	return rsp, nil
}

func ValidateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.EmailId); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	if err := val.SecretCode(req.SecretCode); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
