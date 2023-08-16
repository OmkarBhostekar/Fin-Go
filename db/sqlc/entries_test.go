package simplebank

import (
	"context"
	"testing"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntryById(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	deleteRandomEntry(t, entry1.ID)
}

func TestGetAllEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomEntry(t)
	}
	arg := GetAllEntriesParams{
		Limit:  5,
		Offset: 0,
	}
	entries, err := testQueries.GetAllEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		deleteRandomEntry(t, entry.ID)
	}

}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	arg := UpdateEntryAmountParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}
	err := testQueries.UpdateEntryAmount(context.Background(), arg)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntryById(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 0)
	deleteRandomEntry(t, entry1.ID)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	deleteRandomEntry(t, entry1.ID)
	entry2, err := testQueries.GetEntryById(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)
}

func deleteRandomEntry(t *testing.T, id int64) {
	err := testQueries.DeleteEntryById(context.Background(), id)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntryById(context.Background(), id)
	require.Error(t, err)
	require.Empty(t, entry2)
}

func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: getAccountId(t),
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func getAccountId(t *testing.T) int64 {
	account := CreateRandomAccount(t)
	return account.ID
}
