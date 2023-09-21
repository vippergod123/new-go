package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vipeergod123/simple_bank/util"
)

func createUserTesting(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: util.RandomString(10),
		Email:          util.RandomEmail(),
		FullName:       util.RandomName(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}
func TestCreateUser(t *testing.T) {
	createUserTesting(t)
}

func TestGetUser(t *testing.T) {
	expected := createUserTesting(t)
	actual, err := testQueries.GetUser(context.Background(), expected.Username)

	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.Equal(t, actual.Username, expected.Username)
	require.Equal(t, actual.HashedPassword, expected.HashedPassword)
	require.Equal(t, actual.Email, expected.Email)
	require.Equal(t, actual.FullName, expected.FullName)

	require.NotZero(t, actual.CreatedAt)
	require.True(t, actual.PasswordChangedAt.IsZero())
}

func TestUpdateUser(t *testing.T) {
	account := createAccountTesting(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance,
	}

	expected, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, expected)

	require.Equal(t, account.ID, expected.ID)
	require.Equal(t, account.Balance, expected.Balance)
	require.Equal(t, account.Currency, expected.Currency)
	require.Equal(t, arg.Balance, expected.Balance)
	require.WithinDuration(t, account.CreatedAt, expected.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	actual := createAccountTesting(t)
	err := testQueries.DeleteAccount(context.Background(), actual.ID)
	require.NoError(t, err)

	removeAcc, err := testQueries.GetAccount(context.Background(), actual.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, removeAcc)
}

func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createAccountTesting((t))
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 0)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
