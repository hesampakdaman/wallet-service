package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/hesampakdaman/wallet-service/internal/core/models"
	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

func (s Service) CreateWallet(
	ctx context.Context,
	cmd commands.CreateWallet,
) (models.Wallet, error) {
	logger := s.logger.With(
		slog.String("wallet_id", cmd.WalletID.String()),
		slog.String("player_id", cmd.PlayerID.String()),
		slog.Uint64("initial_balance", uint64(cmd.Balance)),
	)
	logger.DebugContext(ctx, "create wallet: begin.")

	repo, err := s.starter.Repository(ctx, pgx.TxOptions{})
	if err != nil {
		logger.ErrorContext(ctx, "create wallet: tx failed to start", slog.Any("error", err))
		return models.Wallet{}, fmt.Errorf("start tx: %w", err)
	}
	defer repo.Rollback(ctx)

	wallet := models.Wallet{
		ID:       cmd.WalletID,
		PlayerID: cmd.PlayerID,
		Balance:  cmd.Balance,
	}

	if err := repo.Save(ctx, wallet); err != nil {
		logger.ErrorContext(ctx, "create wallet: save failed", slog.Any("error", err))
		return models.Wallet{}, fmt.Errorf("save: %w", err)
	}

	if err := repo.Commit(ctx); err != nil {
		logger.ErrorContext(ctx, "create wallet: commit failed", slog.Any("error", err))
		return models.Wallet{}, fmt.Errorf("commit database: %w", err)
	}

	logger.DebugContext(ctx, "create wallet: complete.")

	return wallet, nil
}
