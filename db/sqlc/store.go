package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	txKey = struct{}{}
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, args CreateUserTxParams) (CreateUserTxResult, error)
	VeirfyEmailTx(ctx context.Context, args VeirfyEmailTxParams) (VeirfyEmailTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
