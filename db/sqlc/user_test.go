package db

import (
	"context"
	"testing"
	"time"

	"github.com/bensmile/wekamakuta/util"
	"github.com/jackc/pgx/v5/pgtype"
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
	user, err := testStore.CreateUser(context.Background(), args)

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

	user, err := testStore.GetUser(context.Background(), user1.Username)

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

func TestUpdateUserOnlyFullName(t *testing.T) {

	oldUser := createRandomUser(t)

	newFullName := util.RandomOwnerName()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, newFullName, oldUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserOnlyEmail(t *testing.T) {

	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, newEmail, oldUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserOnlyPassowrd(t *testing.T) {
	oldUser := createRandomUser(t)
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, newHashedPassword, oldUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)

}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	newEmail := util.RandomEmail()
	newFullName := util.RandomOwnerName()

	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		}, Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		}, FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.NotEqual(t, newHashedPassword, oldUser.HashedPassword)
	require.NotEqual(t, newEmail, oldUser.Email)
	require.NotEqual(t, newFullName, oldUser.FullName)

	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newFullName, updatedUser.FullName)
}
