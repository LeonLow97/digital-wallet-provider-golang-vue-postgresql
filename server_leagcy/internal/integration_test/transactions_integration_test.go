package integration_test

import (
	"context"
	"testing"

	"github.com/LeonLow97/internal/transactions"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestIntegration_GetUserCountById(t *testing.T) {
	testRepo := transactions.NewRepo(testDB)

	count, err := testRepo.GetUserCountByUserId(context.Background(), 1)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, count, 1)

	count, err = testRepo.GetUserCountByUserId(context.Background(), 999)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, count, 0)
}

func TestIntegration_GetUserIdByMobileNumber(t *testing.T) {
	testRepo := transactions.NewRepo(testDB)

	userId, err := testRepo.GetUserIdByMobileNumber(context.Background(), "+44 7712345678")
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, userId, 4)

	userId, err = testRepo.GetUserIdByMobileNumber(context.Background(), "1")
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, userId, 0)
}
