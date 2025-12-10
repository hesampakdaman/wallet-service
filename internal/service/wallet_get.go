package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/hesampakdaman/wallet-service/internal/core/models"
	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

func (s Service) GetWallet(ctx context.Context, cmd commands.GetWallet) (models.Wallet, error) {
	repo, err := s.starter.Repository(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Wallet{}, fmt.Errorf("start tx: %w", err)
	}
	defer repo.Rollback(ctx)

	wallet, err := repo.Get(ctx, cmd.WalletID)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("get wallet: %w", err)
	}

	return wallet, nil
}
