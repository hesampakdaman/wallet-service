package handler

import (
	"errors"
	"net/http"

	"github.com/hesampakdaman/wallet-service/internal/core/errorx"
)

func mapErrToStatusCode(err error) int {
	switch {
	case errors.Is(err, errorx.ErrNotFound):
		return http.StatusNotFound

	case errors.Is(err, errorx.ErrInsufficientBalance):
		return http.StatusBadRequest

	case errors.Is(err, errorx.ErrConflict):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
