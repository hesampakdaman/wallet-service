package service

import (
	"log/slog"

	"github.com/hesampakdaman/wallet-service/internal/adapters/repository"
)

type Service struct {
	logger  *slog.Logger
	starter *repository.Starter
}

func New(logger *slog.Logger, s *repository.Starter) *Service {
	logger = logger.With("component", "service")
	return &Service{logger: logger, starter: s}
}
