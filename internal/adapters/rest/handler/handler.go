package handler

import (
	"log/slog"

	"github.com/hesampakdaman/wallet-service/internal/service"
)

type Handler struct {
	logger *slog.Logger
	svc    *service.Service
}

func New(logger *slog.Logger, s *service.Service) *Handler {
	logger = logger.With("component", "rest-handler")
	return &Handler{logger: logger, svc: s}
}
