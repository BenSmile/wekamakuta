package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bensmile/wekamakuta/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwnerName(), //randomly generated
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {

	account1 := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, account1.Owner, account.Owner)
	require.Equal(t, account1.Balance, account.Balance)
	require.Equal(t, account1.Currency, account.Currency)
	require.Equal(t, account1.CreatedAt, account.CreatedAt)
	require.WithinDuration(t, account1.CreatedAt, account.CreatedAt, time.Second)
	require.Equal(t, account1.ID, account.ID)

}

func TestUpdateAccount(t *testing.T) {

	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
		Balance: util.RandomMoney(),
		ID:      account1.ID,
	}

	account, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, account1.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, account1.Currency, account.Currency)
	require.Equal(t, account1.CreatedAt, account.CreatedAt)
	require.WithinDuration(t, account1.CreatedAt, account.CreatedAt, time.Second)
	require.Equal(t, account1.ID, account.ID)

}

func TestDeleteAccount(t *testing.T) {

	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)

}

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	limit := util.RandomInt(1, 10)

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  int32(limit),
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, int(limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
