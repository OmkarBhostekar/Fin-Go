package simplebank

import (
	"context"
	"testing"
	"time"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	// test bcrypt password matching
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotZero(t, user.CreatedAt)
	return user
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, user1.Username)
	require.Equal(t, user.FullName, user1.FullName)
	require.Equal(t, user.Email, user1.Email)
	require.WithinDuration(t, user.CreatedAt, user1.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, user1.PasswordChangedAt, time.Second)
}
