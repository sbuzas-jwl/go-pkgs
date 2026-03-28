package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/sbuzas-jwl/go-pkgs/todo"
	"github.com/sbuzas-jwl/go-pkgs/todo/http"
	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/logging"
	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/sqlite"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Initialize logging.
	logger := logging.NewLogger("DEBUG")
	todo.SetLogger(logger)
	//Initialize Error Handler
	todo.SetErrorHandler(LogErrorHandler{logger})

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewMain()
	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		todo.Handle(err)
		os.Exit(1)
	}

	logger.Debug("server running", "url", m.HTTPServer.URL())
	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program.
type Main struct {
	// Configuration path and parsed config data.
	Config     Config
	ConfigPath string

	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	HTTPServer *http.Server

	// Services exposed for end-to-end tests.
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Config:     DefaultConfig(),
		ConfigPath: DefaultConfigPath,

		HTTPServer: http.NewServer(),
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	return nil
}

// Run executes the program. The configuration should already be set up before
// calling this function.
func (m *Main) Run(ctx context.Context) (err error) {
	// Instantiate services.
	err = envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   &m.Config,
		Lookuper: envconfig.PrefixLookuper("TODO", envconfig.OsLookuper()),
	})
	if err != nil {
		return fmt.Errorf("unable to process config: %w", err)
	}

	// Copy configuration settings to the HTTP server.
	m.HTTPServer.Port = m.Config.HTTP.Port
	m.HTTPServer.Domain = m.Config.HTTP.Domain

	// Attach underlying services to the HTTP server.
	m.HTTPServer.DB, err = sqlite.NewFromEnv(ctx, &m.Config.DB)
	if err != nil {
		return fmt.Errorf("no database: %w", err)
	}
	// Start the HTTP server.
	if err := m.HTTPServer.Run(ctx); err != nil {
		return err
	}

	// If TLS enabled, redirect non-TLS connections to TLS.
	if m.HTTPServer.UseTLS() {
		go func() {
			log.Fatal(http.ListenAndServeTLSRedirect(m.Config.HTTP.Domain))
		}()
	}

	// Enable internal debug endpoints.
	go func() { http.ListenAndServeDebug() }()

	return nil
}

const (
	// DefaultConfigPath is the default path to the application configuration.
	DefaultConfigPath = "/opt/todo/todod.conf"
)

type Config struct {
	HTTP struct {
		Port   int    `env:"HTTP_PORT"`
		Domain string `env:"HTTP_DOMAIN"`
	}
	DB sqlite.Config
}

// DefaultConfig returns a new instance of Config with defaults set.
func DefaultConfig() Config {
	var config Config
	config.HTTP.Port = 8080
	return config
}

type LogErrorHandler struct {
	log *slog.Logger
}

func (eh LogErrorHandler) Handle(err error) {
	eh.log.With("error", err).Error("unexpected error encountered.")
}
