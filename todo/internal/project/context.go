package project

import (
	"context"
	"log/slog"
	"testing"

	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/logging"
)

func TestContext(t testing.TB) context.Context {
	ctx := context.Background()
	ctx = logging.WithLogger(ctx, TestLogger(t))
	return ctx
}

// TestLogger returns a logger configured for test.
func TestLogger(tb testing.TB) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(tb.Output(), &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
}
