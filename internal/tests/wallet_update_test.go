package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/wallet-service/internal/adapters/rest/handler"
	"github.com/hesampakdaman/wallet-service/internal/core/models"
)

func TestWithdrawShouldReturn200(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	intialBalance := 100
	wallet := createRandomWallet(t, fx, uint(intialBalance))
	expectedBalance := uint(50)

	// Act
	req := handler.UpdateWalletRequest{
		Amount: 50,
		Type:   handler.Withdraw,
	}
	resp, err := fx.DoJSON(http.MethodPatch, fmt.Sprintf("/wallet/%s", wallet.ID), req, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBalance, getWalletBalance(t, fx, wallet.ID))
}

func TestWithdrawInsufficientBalance(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	intialBalance := 100
	wallet := createRandomWallet(t, fx, uint(intialBalance))

	// Act
	req := handler.UpdateWalletRequest{
		Amount: 200,
		Type:   handler.Withdraw,
	}
	resp, err := fx.DoJSON(http.MethodPatch, fmt.Sprintf("/wallet/%s", wallet.ID), req, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDepositShouldReturn200(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	intialBalance := 100
	wallet := createRandomWallet(t, fx, uint(intialBalance))
	expectedBalance := uint(150)

	// Act
	req := handler.UpdateWalletRequest{
		Amount: 50,
		Type:   handler.Deposit,
	}
	resp, err := fx.DoJSON(http.MethodPatch, fmt.Sprintf("/wallet/%s", wallet.ID), req, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedBalance, getWalletBalance(t, fx, wallet.ID))
}

func TestNonExistingWalletUpdateShouldReturn404(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)

	// Act
	req := handler.UpdateWalletRequest{
		Amount: 50,
		Type:   handler.Deposit,
	}
	resp, err := fx.DoJSON(http.MethodPatch, fmt.Sprintf("/wallet/%s", uuid.NewString()), req, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNonExistingTransactionTypeShouldReturn400(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)

	// Act
	req := handler.UpdateWalletRequest{
		Amount: 50,
		Type:   "random",
	}
	resp, err := fx.DoJSON(http.MethodPatch, fmt.Sprintf("/wallet/%s", uuid.NewString()), req, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func createRandomWallet(t *testing.T, fx Fixture, initialBalance uint) models.Wallet {
	t.Helper()
	playerID := uuid.New()
	req := handler.CreateWalletRequest{
		PlayerID:       playerID,
		InitialBalance: initialBalance,
	}

	var payload handler.CreateWalletResponse
	resp, err := fx.DoJSON(http.MethodPost, "/wallet", req, &payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return models.Wallet{
		ID:       payload.WalletID,
		PlayerID: playerID,
		Balance:  initialBalance,
	}
}

func getWalletBalance(t *testing.T, fx Fixture, walletID uuid.UUID) uint {
	t.Helper()

	var payload handler.GetWalletResponse
	resp, err := fx.DoJSON(
		http.MethodGet,
		fmt.Sprintf("/wallet/%s", walletID.String()),
		nil,
		&payload,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	return payload.Balance
}
