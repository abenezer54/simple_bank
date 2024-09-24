package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/abenezer54/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	res, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, args.Owner, res.Owner)
	require.Equal(t, args.Balance, res.Balance)
	require.Equal(t, args.Currency, res.Currency)

	require.NotZero(t, res.ID)
	require.NotZero(t, res.CreatedAt)
	return res
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	res, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, res.ID, account.ID)
	require.Equal(t, res.Owner, account.Owner)
	require.Equal(t, res.Balance, account.Balance)
	require.Equal(t, res.Currency, account.Currency)
	require.Equal(t, res.CreatedAt, account.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	res, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, res.ID, account.ID)
	require.Equal(t, res.Owner, account.Owner)
	require.Equal(t, res.Currency, account.Currency)
	require.Equal(t, res.CreatedAt, account.CreatedAt)

	require.Equal(t, res.Balance, args.Balance) // check the balance with the new arg.balance
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	res, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, res)
}


func TestListAccounts(t *testing.T) {
	for i:=1; i<=10; i++{
		createRandomAccount(t)
	}
	args := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}
	res, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Len(t, res, 5)
	for _, account := range res{
		require.NotEmpty(t, account)
	}

}
