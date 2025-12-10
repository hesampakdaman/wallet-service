package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

type TransactionType string

var (
	Withdraw TransactionType = "withdraw"
	Deposit  TransactionType = "deposit"
)

type UpdateWalletRequest struct {
	Amount uint            `json:"amount"`
	Type   TransactionType `json:"type"`
}

func (h Handler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.PathValue("id") == "" {
		http.Error(w, "missing wallet id", http.StatusBadRequest)
		return
	}

	walletID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "wallet id must be UUID", http.StatusBadRequest)
		return
	}

	var request UpdateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "could not parse argument", http.StatusBadRequest)
		return
	}

	var updateErr error
	switch request.Type {
	case Withdraw:
		updateErr = h.svc.Withdraw(
			r.Context(),
			commands.Withdraw{WalletID: walletID, Amount: request.Amount},
		)
	case Deposit:
		updateErr = h.svc.Deposit(
			r.Context(),
			commands.Deposit{WalletID: walletID, Amount: request.Amount},
		)
	default:
		http.Error(w, "transaction type must be withdraw or deposit", http.StatusBadRequest)
		return
	}

	if updateErr != nil {
		http.Error(w, "could not update wallet", mapErrToStatusCode(updateErr))
		return
	}
}
