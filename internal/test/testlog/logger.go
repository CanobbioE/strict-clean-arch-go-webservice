package testlog

import (
	"context"
	"log/slog"
)

type discardHandler struct{}

func (discardHandler) Enabled(_ context.Context, _ slog.Level) bool { return false }

//nolint:gocritic // necessary to satisfy interface
func (discardHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (discardHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return discardHandler{} }
func (discardHandler) WithGroup(_ string) slog.Handler               { return discardHandler{} }

// NewTestLogger returns a noop logger for unit tests.
func NewTestLogger() *slog.Logger {
	return slog.New(discardHandler{})
}
