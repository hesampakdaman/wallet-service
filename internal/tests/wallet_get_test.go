package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/wallet-service/internal/adapters/rest/handler"
)

func TestGetWalletShouldReturn200(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)

	var createPayload handler.CreateWalletResponse
	resp, err := fx.DoJSON(
		http.MethodPost,
		"/wallet",
		handler.CreateWalletRequest{
			PlayerID:       uuid.New(),
			InitialBalance: 100,
		},
		&createPayload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Act
	var payload handler.GetWalletResponse
	resp, err = fx.DoJSON(
		http.MethodGet,
		fmt.Sprintf("/wallet/%s", createPayload.WalletID),
		nil,
		&payload,
	)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, uint(100), payload.Balance)
}

func TestGetWalletShouldReturn404(t *testing.T) {
	t.Parallel()
	fx := NewFixture(t)

	// Arrange
	randomID := uuid.NewString()

	// Act
	resp, err := fx.DoJSON(http.MethodGet, fmt.Sprintf("/wallet/%s", randomID), nil, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
