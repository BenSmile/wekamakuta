package gapi

import (
	"fmt"

	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/pb"
	"github.com/bensmile/wekamakuta/token"
	"github.com/bensmile/wekamakuta/util"
	"github.com/bensmile/wekamakuta/worker"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer created a new gRPC srver and setup routing
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
