package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

func (s Service) Deposit(ctx context.Context, cmd commands.Deposit) error {
	logger := s.logger.With(
		slog.String("wallet_id", cmd.WalletID.String()),
		slog.Uint64("amount", uint64(cmd.Amount)),
	)
	logger.DebugContext(ctx, "deposit: begin.")

	repo, err := s.starter.Repository(ctx, pgx.TxOptions{})
	if err != nil {
		logger.ErrorContext(ctx, "deposit: tx failed to start", slog.Any("error", err))
		return fmt.Errorf("start tx: %w", err)
	}
	defer repo.Rollback(ctx)

	if err := repo.AdjustBalance(ctx, cmd.WalletID, int(cmd.Amount)); err != nil {
		logger.ErrorContext(ctx, "deposit: save failed", slog.Any("error", err))
		return fmt.Errorf("save: %w", err)
	}

	if err := repo.Commit(ctx); err != nil {
		logger.ErrorContext(ctx, "deposit: commit failed", slog.Any("error", err))
		return fmt.Errorf("commit database: %w", err)
	}

	logger.DebugContext(ctx, "amount deposited.")

	return nil
}
