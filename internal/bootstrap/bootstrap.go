package bootstrap

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hesampakdaman/wallet-service/internal/adapters/repository"
	"github.com/hesampakdaman/wallet-service/internal/adapters/rest"
	"github.com/hesampakdaman/wallet-service/internal/service"
)

type App struct {
	Service *service.Service
	Router  http.Handler
}

func NewApp(logger *slog.Logger, pool *pgxpool.Pool) App {
	st := repository.NewStarter(pool)
	svc := service.New(logger, st)
	router := rest.NewRouter(logger, svc)
	return App{Service: svc, Router: router}
}
