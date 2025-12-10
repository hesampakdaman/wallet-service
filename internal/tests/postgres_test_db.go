package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/wallet-service/internal/adapters/repository/migrations"
)

type PostgresTestDB struct {
	adminPool  *pgxpool.Pool
	connString string
	dbName     string
}

func NewPostgresTestDB(ctx context.Context, conn, dbName string) (*PostgresTestDB, error) {
	p, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}

	return &PostgresTestDB{
		adminPool:  p,
		connString: conn,
		dbName:     dbName,
	}, nil
}

func (p PostgresTestDB) NewDBConnString(newName string) string {
	return strings.Replace(p.connString, p.dbName, newName, 1)
}

func (p PostgresTestDB) TestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := t.Context()
	dbName := "test-db_" + uuid.New().String()

	_, err := p.adminPool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	require.NoError(t, err)

	cfg, err := pgxpool.ParseConfig(p.NewDBConnString(dbName))
	require.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	p.Migrate(t, cfg.ConnString())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}

func (p PostgresTestDB) Migrate(t *testing.T, conn string) {
	src, err := iofs.New(migrations.Files, ".")
	if err != nil {
		t.Fatalf("iofs new: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, conn)
	if err != nil {
		t.Fatalf("migrate new: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("migrate up: %v", err)
	}
}
