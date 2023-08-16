package simplebank

import (
	"context"
	"testing"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransferById(t *testing.T) {
	transfer1, from, to := CreateRandomTransfer(t)
	_ = from
	_ = to
	transfer2, err := testQueries.GetTransferById(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, 0)
}

func TestGetTransfersByFromAccountId(t *testing.T) {
	transfer1, from, to := CreateRandomTransfer(t)
	_ = from
	_ = to
	transfer2, err := testQueries.GetTransfersByFromAccountId(context.Background(), transfer1.FromAccountID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2[0].ID)
	require.Equal(t, transfer1.FromAccountID, transfer2[0].FromAccountID)
}

func TestGetTransfersByToAccountId(t *testing.T) {
	transfer1, from, to := CreateRandomTransfer(t)
	_ = from
	_ = to
	transfer2, err := testQueries.GetTransfersByToAccountId(context.Background(), transfer1.ToAccountID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2[0].ID)
	require.Equal(t, transfer1.ToAccountID, transfer2[0].ToAccountID)
}

func TestGetTransfersByFromAccountIdAndToAccountId(t *testing.T) {
	transfer1, from, to := CreateRandomTransfer(t)
	_ = from
	_ = to
	args := GetTransfersByFromAccountIdAndToAccountIdParams{
		FromAccountID: transfer1.FromAccountID,
		ToAccountID:   transfer1.ToAccountID,
	}
	transfer2, err := testQueries.GetTransfersByFromAccountIdAndToAccountId(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2[0].ID)
	require.Equal(t, transfer1.FromAccountID, transfer2[0].FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2[0].ToAccountID)
}

func TestDeleteTransferById(t *testing.T) {
	transfer1, from, to := CreateRandomTransfer(t)
	_ = from
	_ = to
	err := testQueries.DeleteTransferById(context.Background(), transfer1.ID)
	require.NoError(t, err)
	transfer2, err := testQueries.GetTransferById(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Empty(t, transfer2)
}

func TestGetAllTransfers(t *testing.T) {
	CreateRandomTransfer(t)
}

func CreateRandomTransfer(t *testing.T) (Transfer, Account, Account) {

	acTo := CreateRandomAccount(t)
	acFrom := CreateRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: acTo.ID,
		ToAccountID:   acFrom.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer, acTo, acFrom
}
