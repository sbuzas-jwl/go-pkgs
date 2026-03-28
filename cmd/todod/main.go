package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sbuzas-jwl/go-pkgs/internal/logging"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Initialize logging.
	logger := logging.NewLogger("DEBUG")
	ctx = logging.WithLogger(ctx, logger)
	// Instantiate a new type to represent our application.
	m := NewMain(ctx)
	if err := m.Open(); err != nil {
		m.Close()
		os.Exit(1)
	}

	logger.Debug("server running")
	// Wait for CTRL-C.
	<-ctx.Done()

	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
