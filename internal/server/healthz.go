package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sbuzas-jwl/go-pkgs/internal/cache"
	"github.com/sbuzas-jwl/go-pkgs/internal/logging"
	"github.com/sbuzas-jwl/go-pkgs/internal/sqlite"
)

func HandleHealthz(db *sqlite.DB) http.Handler {

	cacher, _ := cache.New[bool](1 * time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx).WithGroup("server.healthz")

		result, _ := cacher.WriteThruLookup("healthz", func() (bool, error) {
			ctx, cancel := context.WithTimeoutCause(ctx, 2*time.Second,
				errors.New("healthcheck timeout exceeded"))
			defer cancel()
			conn, err := db.Pool.Connx(ctx)
			if err != nil {
				logger.Error("failed to acquire database connection", "error", err)
				return false, nil
			}
			defer conn.Close() // nolint:errcheck

			if err := conn.PingContext(ctx); err != nil {
				logger.Error("failed to ping database", "error", err)
				return false, nil
			}

			return true, nil
		})

		if !result {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok"}`)
	})
}
