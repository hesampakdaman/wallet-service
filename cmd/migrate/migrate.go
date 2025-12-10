package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/hesampakdaman/wallet-service/internal/adapters/repository/migrations"
)

func main() {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL env variable must be set")
	}

	src, err := iofs.New(migrations.Files, ".")
	if err != nil {
		log.Fatalf("iofs new: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("migrate new: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
}
