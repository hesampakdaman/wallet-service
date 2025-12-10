package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

type GetWalletResponse struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Balance  uint      `json:"balance"`
}

func (h Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		http.Error(w, "missing wallet id", http.StatusBadRequest)
		return
	}

	walletID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "wallet id must be UUID", http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.GetWallet(r.Context(), commands.GetWallet{WalletID: walletID})
	if err != nil {
		http.Error(w, "could not get wallet", mapErrToStatusCode(err))
		return
	}

	if err := json.NewEncoder(w).Encode(GetWalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	}); err != nil {
		http.Error(w, "could not encode json response", http.StatusInternalServerError)
		return
	}
}
