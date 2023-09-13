package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vipeergod123/simple_bank/util"
)

func createAccountTesting(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createAccountTesting(t)
}

func TestGetAccount(t *testing.T) {
	actual := createAccountTesting(t)
	expected, err := testQueries.GetAccount(context.Background(), actual.ID)

	require.NoError(t, err)
	require.NotEmpty(t, expected)

	require.Equal(t, expected.Owner, actual.Owner)
	require.Equal(t, expected.Balance, actual.Balance)
	require.Equal(t, expected.Currency, actual.Currency)

	require.NotZero(t, expected.ID)
	require.NotZero(t, expected.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
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

func TestDeleteAccount(t *testing.T) {
	actual := createAccountTesting(t)
	err := testQueries.DeleteAccount(context.Background(), actual.ID)
	require.NoError(t, err)

	removeAcc, err := testQueries.GetAccount(context.Background(), actual.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, removeAcc)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createAccountTesting((t))
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
