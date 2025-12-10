package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *PostgresTestDB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, conn, dbName := initPostgresContainer(ctx)
	db = must(NewPostgresTestDB(ctx, conn, dbName))

	code := m.Run()

	go func() { _ = testcontainers.TerminateContainer(pgContainer) }()

	os.Exit(code)
}

func initPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, string, string) {
	dbName := "test-db"
	container := must(postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithDatabase(dbName),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		postgres.WithSQLDriver("pgx"),
	))
	connString := must(container.ConnectionString(ctx, "sslmode=disable"))
	return container, connString, dbName
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	return v
}
