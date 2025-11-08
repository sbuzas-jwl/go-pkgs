package todo

import (
	"log/slog"

	"github.com/sbuzas-jwl/go-pkgs/todo/internal/global"
)

// SetLogger configures the logger used internally.
func SetLogger(logger *slog.Logger) {
	global.SetLogger(logger)
}
