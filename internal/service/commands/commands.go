package commands

import (
	"github.com/google/uuid"
)

type CreateWallet struct {
	WalletID uuid.UUID
	PlayerID uuid.UUID
	Balance  uint
}

type GetWallet struct {
	WalletID uuid.UUID
}

type Deposit struct {
	WalletID uuid.UUID
	Amount   uint
}

type Withdraw struct {
	WalletID uuid.UUID
	Amount   uint
}
