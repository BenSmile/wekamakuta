package db

import (
	"context"
	"testing"
	"time"

	"github.com/bensmile/wekamakuta/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	args := CreateUserParams{
		Username:       util.RandomOwnerName(), //randomly generated
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwnerName(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, user1.Username, user.Username)
	require.Equal(t, user1.FullName, user.FullName)
	require.Equal(t, user1.HashedPassword, user.HashedPassword)
	require.Equal(t, user1.CreatedAt, user.CreatedAt)
	require.WithinDuration(t, user1.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user.PasswordChangedAt, time.Second)

}
