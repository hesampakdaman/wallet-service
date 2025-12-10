package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hesampakdaman/wallet-service/internal/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))
	logger.InfoContext(ctx, "Starting service on 8080...")

	pool := createPool(ctx)
	app := bootstrap.NewApp(logger, pool)
	server := http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.ErrorContext(ctx, "Failed to listen and serve HTTP.")
			stop()
		}
	}()

	logger.InfoContext(ctx, "Service started.")
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
	pool.Close()
}

func createPool(ctx context.Context) *pgxpool.Pool {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL env variable must be set")
	}

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
