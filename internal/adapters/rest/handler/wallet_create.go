package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

type CreateWalletRequest struct {
	PlayerID       uuid.UUID `json:"player_id"`
	InitialBalance uint      `json:"initial_balance"`
}

type CreateWalletResponse struct {
	WalletID uuid.UUID `json:"wallet_id"`
}

func (h Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	var request CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "could not parse argument", http.StatusBadRequest)
		return
	}

	if request.PlayerID == uuid.Nil {
		http.Error(w, "player id must be non-nil UUID", http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.CreateWallet(ctx, commands.CreateWallet{
		WalletID: uuid.New(),
		PlayerID: request.PlayerID,
		Balance:  request.InitialBalance,
	})
	if err != nil {
		http.Error(w, "could not create wallet", mapErrToStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(CreateWalletResponse{WalletID: wallet.ID}); err != nil {
		h.logger.ErrorContext(
			ctx,
			"could not write wallet id to http output",
			slog.String("error", err.Error()),
		)
		http.Error(w, "could not write output data", http.StatusInternalServerError)
		return
	}
}
