package global

import (
	"log/slog"
	"os"
	"sync/atomic"
)

// The default logger uses slog which is backed by the standard `log.Logger`
// interface. This logger will only show messages at the Error Level.
var globalLogger = func() *atomic.Pointer[slog.Logger] {
	l := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	p := new(atomic.Pointer[slog.Logger])
	p.Store(l)
	return p
}()

func SetLogger(l *slog.Logger) {
	globalLogger.Store(l)
}

// GetLogger returns the global logger.
func GetLogger() *slog.Logger {
	return globalLogger.Load()
}

// Info prints messages about the general state of the API or SDK.
// This should usually be less than 5 messages a minute.
func Info(msg string, keysAndValues ...interface{}) {
	GetLogger().Info(msg, keysAndValues...)
}

// Error prints messages about exceptional states of the API or SDK.
func Error(err error, msg string, keysAndValues ...interface{}) {
	GetLogger().
		With("error", slog.AnyValue(err)).
		Error(msg, keysAndValues...)
}

// Debug prints messages about all internal changes in the API or SDK.
func Debug(msg string, keysAndValues ...interface{}) {
	GetLogger().Debug(msg, keysAndValues...)
}

// Warn prints messages about warnings in the API or SDK.
// Not an error but is likely more important than an informational event.
func Warn(msg string, keysAndValues ...interface{}) {
	GetLogger().Warn(msg, keysAndValues...)
}
