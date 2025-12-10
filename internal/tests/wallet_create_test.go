package tests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/wallet-service/internal/adapters/rest/handler"
)

func TestCreateWalletShouldReturn201(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	req := handler.CreateWalletRequest{
		PlayerID:       uuid.New(),
		InitialBalance: 100,
	}

	// Act
	var payload handler.CreateWalletResponse
	resp, err := fx.DoJSON(http.MethodPost, "/wallet", req, &payload)
	require.NoError(t, err)

	// Assert
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.Nil, payload.WalletID)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateWalletInvalidPlayerID(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	req := map[string]any{
		"player_id":       "123",
		"initial_balance": 100,
	}

	// Act
	resp, _, err := fx.DoJSONRaw(http.MethodPost, "/wallet", req)
	require.NoError(t, err)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateWalletNegativeIntialBalance(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	req := map[string]any{
		"player_id":       uuid.NewString(),
		"initial_balance": -100,
	}

	// Act
	resp, _, err := fx.DoJSONRaw(http.MethodPost, "/wallet", req)
	require.NoError(t, err)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateWalletDuplicatePlayerIDReturns409(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	playerID := uuid.New()

	firstReq := handler.CreateWalletRequest{
		PlayerID:       playerID,
		InitialBalance: 100,
	}
	secondReq := handler.CreateWalletRequest{
		PlayerID:       playerID,
		InitialBalance: 200,
	}

	// Act
	resp, err := fx.DoJSON(http.MethodPost, "/wallet", firstReq, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = fx.DoJSON(http.MethodPost, "/wallet", secondReq, nil)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}
