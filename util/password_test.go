package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestCheckPassword(t *testing.T) {
	hash1, err := HashPassword("password")
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	err = CheckPassword("password", hash1)
	require.NoError(t, err)

	err = CheckPassword("wrong", hash1)
	require.Error(t, err)

}
