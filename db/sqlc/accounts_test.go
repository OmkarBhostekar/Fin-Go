package simplebank

import (
	"context"
	"testing"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreate5Account(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomAccount(t)
	}
}
func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
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

func TestGetAccountById(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account, err := testQueries.GetAcountById(context.Background(), int64(account1.ID))
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account1.ID, account.ID)
	require.Equal(t, account1.Owner, account.Owner)
	require.Equal(t, account1.Balance, account.Balance)
	require.Equal(t, account1.Currency, account.Currency)
	require.NotZero(t, account.ID)
}

func TestGetAllAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 5; i++ {
		lastAccount = CreateRandomAccount(t)
	}
	arg := GetAllAccountsParams{
		Limit:  5,
		Offset: 0,
		Owner:  lastAccount.Owner,
	}
	accounts, err := testQueries.GetAllAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestDeleteAccountById(t *testing.T) {
	account1 := CreateRandomAccount(t)
	id := account1.ID
	err := testQueries.DeleteAccountById(context.Background(), int64(id))
	require.NoError(t, err)
	account, err := testQueries.GetAcountById(context.Background(), int64(id))
	require.Error(t, err)
	require.Empty(t, account)
}
