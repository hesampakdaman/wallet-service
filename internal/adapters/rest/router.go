package rest

import (
	"log/slog"
	"net/http"

	"github.com/hesampakdaman/wallet-service/internal/adapters/rest/handler"
	"github.com/hesampakdaman/wallet-service/internal/service"
)

func NewRouter(logger *slog.Logger, svc *service.Service) *http.ServeMux {
	handler := handler.New(logger, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /wallet/{id}", handler.GetWallet)
	mux.HandleFunc("PATCH /wallet/{id}", handler.UpdateWallet)
	mux.HandleFunc("POST /wallet", handler.CreateWallet)

	return mux
}
