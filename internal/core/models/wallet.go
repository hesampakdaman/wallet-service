package models

import (
	"github.com/google/uuid"

	"github.com/hesampakdaman/wallet-service/internal/core/errorx"
)

type Wallet struct {
	ID       uuid.UUID
	PlayerID uuid.UUID
	Balance  uint
}

func (w *Wallet) Add(amount uint) {
	w.Balance += amount
}

func (w *Wallet) Sub(amount uint) error {
	if amount > w.Balance {
		return errorx.ErrInsufficientBalance
	}

	w.Balance -= amount

	return nil
}
