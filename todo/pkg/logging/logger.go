// Package logging sets up and configures logging.
package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logger is stored.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

// ParseLevel converts a string representation of a log level to a slog.Level.
func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	if err != nil {
		return 0, fmt.Errorf("invalid log level string: %w", err)
	}
	return level, nil
}

// NewLogger creates a new logger with the given configuration.
func NewLogger(level string) *slog.Logger {
	slogLevel, err := ParseLevel(level)
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		}),
	)
	if err != nil {
		logger.Info("unable to parse log level, proceeding with INFO",
			slog.String("level", level))
	}

	return logger
}

// NewLoggerFromEnv creates a new logger from the environment. It consumes
// LOG_LEVEL for determining the level.
func NewLoggerFromEnv() *slog.Logger {
	level := os.Getenv("LOG_LEVEL")
	return NewLogger(level)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}
