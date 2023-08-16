package simplebank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Omkar",
		Balance:  2000,
		Currency: "INR",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}

func TestGetAccountById(t *testing.T) {
	id := 1
	account, err := testQueries.GetAcountById(context.Background(), int64(id))
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, int64(id), account.ID)
	require.NotZero(t, account.ID)
}

func TestGetAllAccounts(t *testing.T) {
	arg := GetAllAccountsParams{
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.GetAllAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
}
