package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/hesampakdaman/wallet-service/internal/core/errorx"
	"github.com/hesampakdaman/wallet-service/internal/service/commands"
)

func (s Service) Withdraw(ctx context.Context, cmd commands.Withdraw) error {
	logger := s.logger.With(
		slog.String("wallet_id", cmd.WalletID.String()),
		slog.Uint64("amount", uint64(cmd.Amount)),
	)
	logger.DebugContext(ctx, "withdraw: begin.")

	repo, err := s.starter.Repository(ctx, pgx.TxOptions{})
	if err != nil {
		logger.ErrorContext(ctx, "withdraw: tx failed to start", slog.Any("error", err))
		return fmt.Errorf("withdraw: start tx: %w", err)
	}
	defer repo.Rollback(ctx)

	wallet, err := repo.GetForUpdate(ctx, cmd.WalletID)
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			logger.DebugContext(ctx, "withdraw: wallet not found")
			return err
		}
		return fmt.Errorf("get wallet for update: %w", err)
	}

	if err := wallet.Sub(cmd.Amount); err != nil {
		logger.DebugContext(ctx, "withdraw: insufficient funds")
		return err
	}

	if err := repo.Save(ctx, wallet); err != nil {
		logger.ErrorContext(ctx, "withdraw: save failed", slog.Any("error", err))
		return fmt.Errorf("save: %w", err)
	}

	if err := repo.Commit(ctx); err != nil {
		logger.ErrorContext(ctx, "deposit: commit failed", slog.Any("error", err))
		return fmt.Errorf("commit database: %w", err)
	}

	logger.DebugContext(ctx, "withdraw: complete.")

	return nil
}
