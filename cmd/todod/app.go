package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sbuzas-jwl/go-pkgs/internal/logging"
	"github.com/sbuzas-jwl/go-pkgs/internal/server"
	"github.com/sethvargo/go-envconfig"
)

type HTTPServer interface {
	http.Handler
	io.Closer
	Open() error

	UseTLS() bool
}

// Main represents the program.
type Main struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	log    *slog.Logger
	srv    *http.Server

	// Configuration path and parsed config data.
	Config     Config
	ConfigPath string

	// HTTP server for handling HTTP communication.
	HTTPServer HTTPServer

	// Services exposed for end-to-end tests.
}

// NewMain returns a new instance of Main.
//
// This type lets us share setup code with our end-to-end tests.
func NewMain(ctx context.Context) *Main {
	ctx, cancel := context.WithCancelCause(context.WithoutCancel(ctx))
	return &Main{
		ctx:        ctx,
		cancel:     cancel,
		log:        logging.FromContext(ctx),
		Config:     DefaultConfig(),
		ConfigPath: DefaultConfigPath,
		HTTPServer: nil,
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	//TODO: cleanup resources
	if m.srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := m.srv.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.log.Error(err.Error())
		}
	}

	return nil
}

// Run executes the program. The configuration should already be set up before
// calling this function.
func (m *Main) Open() (err error) {
	if err = envconfig.ProcessWith(m.ctx, &envconfig.Config{
		Target:   &m.Config,
		Lookuper: envconfig.PrefixLookuper("TODO", envconfig.OsLookuper()),
	}); err != nil {
		return fmt.Errorf("unable to process config: %w", err)
	}
	// Optionally: Override logging based on ENV
	if level, ok := os.LookupEnv("LOG_LEVEL"); ok {
		m.log = logging.NewLogger(level)
	}
	// TODO: Initialize App HTTPServer
	// todoSrv := todohttp.NewServer()
	// app := todo.New()
	// todoSrv.App = app

	// TODO: Copy configuration settings to the HTTP server.

	// TODO: Attach underlying services to the HTTP server.
	// m.HTTPServer.DB, err = sqlite.NewFromEnv(ctx, &m.Config.DB)
	// if err != nil {
	// 	return fmt.Errorf("no database: %w", err)
	// }

	// go func() {
	// 	err := m.HTTPServer.Open()
	// 	if err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		m.log.Error("HTTPServer terminated unexpectedly",
	// 			slog.String("err.message", err.Error()))
	// 	}
	// }()

	m.log.Debug("Starting web server...")
	handler := server.RequestLogger(m.log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
			logging.FromContext(r.Context()).Error(err.Error())
		}
	}))
	m.srv = &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		if err := m.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			m.log.Error(err.Error())
		} else {
			m.log.Debug("Shutting down web server...")
		}
	}()
	//TODO: implement a secure redirect server in internal
	// If TLS enabled, redirect non-TLS connections to TLS.
	// if m.HTTPServer.UseTLS() {
	// go func() {
	// 	log.Fatal(server.ListenAndServeTLSRedirect(m.Config.HTTP.Domain))
	// }()
	// }

	// Enable internal debug endpoints.
	if m.UseDebug() {
		// TODO: implement a debug server in internal
		// go func() { todohttp.ListenAndServeDebug() }()
	}

	return nil
}

func (m *Main) UseDebug() bool {
	return m.Config.Debug != ""
}
