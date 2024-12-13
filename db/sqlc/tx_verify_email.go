package db

import (
	"context"
	"database/sql"
	"fmt"
)

type VeirfyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VeirfyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VeirfyEmailTx(ctx context.Context, args VeirfyEmailTxParams) (VeirfyEmailTxResult, error) {
	var result VeirfyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		fmt.Printf("%+v", args)
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         args.EmailId,
			SecretCode: args.SecretCode,
		})

		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		})

		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
