package errorx

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrConflict            = errors.New("conflict")
)
