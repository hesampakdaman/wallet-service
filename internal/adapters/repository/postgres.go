package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hesampakdaman/wallet-service/internal/core/errorx"
	"github.com/hesampakdaman/wallet-service/internal/core/models"
)

type Repository struct {
	tx pgx.Tx
}

func (r Repository) Rollback(ctx context.Context) {
	_ = r.tx.Rollback(ctx)
}

func (r Repository) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

func (r Repository) Save(ctx context.Context, w models.Wallet) error {
	_, err := r.tx.Exec(ctx, `
        INSERT INTO wallets (id, player_id, balance)
        VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET balance = EXCLUDED.balance
        `,
		w.ID.String(),
		w.PlayerID.String(),
		w.Balance,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		// Unique constraint on player_id hit.
		if strings.EqualFold(pgErr.ConstraintName, "wallets_player_id_key") {
			return errorx.ErrConflict
		}
	}

	return err
}

func (r Repository) Get(ctx context.Context, walletID uuid.UUID) (models.Wallet, error) {
	var w models.Wallet
	if err := r.tx.QueryRow(ctx, `
        SELECT id, player_id, balance FROM wallets
        WHERE id = $1
        `,
		walletID.String(),
	).Scan(&w.ID, &w.PlayerID, &w.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return w, errorx.ErrNotFound
		}
		return w, err
	}

	return w, nil
}

func (r Repository) AdjustBalance(ctx context.Context, walletID uuid.UUID, amount int) error {
	ct, err := r.tx.Exec(ctx, `
        UPDATE wallets
        SET balance = balance + $1
        WHERE id = $2
	    `,
		amount,
		walletID.String(),
	)
	if err == nil && ct.RowsAffected() == 0 { // If no wallet exists
		return errorx.ErrNotFound
	}
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			if strings.EqualFold(pgErr.ConstraintName, "wallets_player_id_key") {
				return errorx.ErrConflict
			}
		case pgerrcode.CheckViolation:
			if strings.EqualFold(pgErr.ConstraintName, "wallets_balance_non_negative") {
				return errorx.ErrInsufficientBalance
			}
		}
	}

	return err
}
