package gapi

import (
	"context"
	"time"

	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/pb"
	"github.com/bensmile/wekamakuta/util"
	"github.com/bensmile/wekamakuta/val"
	"github.com/bensmile/wekamakuta/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := ValidateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	txResult, err := server.store.CreateUserTx(ctx,
		db.CreateUserTxParams{
			CreateUserParams: db.CreateUserParams{
				Username:       req.GetUsername(),
				Email:          req.GetEmail(),
				HashedPassword: hashedPassword,
				FullName:       req.GetFullName(),
			},
			AfterCreate: func(user db.User) error {
				// send verification email
				// TODO: use db transaction
				opts := []asynq.Option{
					asynq.MaxRetry(10),
					asynq.ProcessIn(10 * time.Second),
					asynq.Queue(worker.QUEUE_CRITICAL),
				}
				return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}, opts...)
			},
		})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}

	return rsp, nil
}

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateFullname(req.FullName); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := val.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}
