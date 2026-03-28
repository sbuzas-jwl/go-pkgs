package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sbuzas-jwl/go-pkgs/todo"
	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/logging"
	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/server"
	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/sqlite"
	"golang.org/x/crypto/acme/autocert"
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 5 * time.Second

type Server struct {
	server *server.Server
	// server *http.Server
	router *mux.Router

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocert.
	Port   int
	Domain string

	// Infra Dependencies
	DB *sqlite.DB

	// Services used by the various HTTP routes.
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	// Create a new server that wraps the net/http server & add a gorilla router.
	s := &Server{
		router: mux.NewRouter(),
	}

	s.router.Use(reportPanic)
	// Our router can be wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes changing route paths for JSON endpoints & overriding methods.

	// Setup error handling routes.
	s.router.NotFoundHandler = server.RequestLogger(logging.DefaultLogger())(
		http.HandlerFunc(s.handleNotFound),
	)
	return s
}

func (s *Server) configureRouter(ctx context.Context) {
	// Setup a base router that for api endpoints.
	router := s.router.PathPrefix("/").Subrouter()
	router.Use(
		s.authenticate,
		server.RequestLogger(logging.FromContext(ctx)),
	)
	// add additional global middlewares here
	{
		r := router.Host("localhost:8080").Subrouter()
		r.Handle("/debug", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.WriteHeader(http.StatusTeapot)
			fmt.Fprintf(w, `{"status":"i am a little teapot"}`)
		}))
	}
	// Register unauthenticated routes.
	{
		r := router.PathPrefix("/").Subrouter()
		r.Handle("/healthz", server.HandleHealthz(s.DB))
	}

	// Register authenticated routes.
	// {
	// 	r := router.PathPrefix("/").Subrouter()
	// 	r.Use(s.requireAuth)
	// 	r.HandleFunc("/settings", s.handleSettings).Methods("GET")
	// 	s.registerDialRoutes(r)
	// 	s.registerDialMembershipRoutes(r)
	// 	s.registerEventRoutes(r)
	// }
}

func (s *Server) Run(ctx context.Context) error {
	// Validate settings.
	s.configureRouter(ctx)
	// Open a listener on our bind address.
	var listner net.Listener
	if s.Domain != "" {
		listner = autocert.NewListener(s.Domain)
	} else {
		addr := fmt.Sprintf(":%d", s.Port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		listner = ln
	}

	server, err := server.NewFromListener(listner)
	if err != nil {
		todo.Handle(err)
	}

	s.server = server
	go func() {
		if err := s.server.ServeHTTPHandler(ctx, s.router); !errors.Is(err, http.ErrServerClosed) {
			todo.Handle(err)
		}
	}()

	return nil
}

// UseTLS returns true if the cert & key file are specified.
func (s *Server) UseTLS() bool {
	return s.Domain != ""
}

// Scheme returns the URL scheme for the server.
func (s *Server) Scheme() string {
	if s.UseTLS() {
		return "https"
	}
	return "http"
}

// URL returns the local base URL of the running server.
func (s *Server) URL() string {
	scheme, port := s.Scheme(), s.server.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.server.Port())
}

// reportPanic is middleware for catching panics and reporting them.
func reportPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if _err, eok := err.(error); eok {
					todo.Handle(_err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// handleNotFound handles requests to routes that don't exist.
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// authenticate is middleware for authenticating a request, and adding authorization context to the request chain.
func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add Authentication
	})
}

// ListenAndServeTLSRedirect runs an HTTP server on port 80 to redirect users
// to the TLS-enabled port 443 server.
func ListenAndServeTLSRedirect(domain string) error {
	return http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+domain, http.StatusFound)
	}))
}

// ListenAndServeDebug runs an HTTP server with /debug endpoints (e.g. pprof, vars).
func ListenAndServeDebug() error {
	h := http.NewServeMux()
	// TODO: metrics exporter endpoint h.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":6060", h)
}
