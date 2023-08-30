// internal/infrastructure/http/http_server.go
package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
)

var logger = logging.NewLogger()

type HttpServer struct{}

// NewHTTPServer creates a new instance of HttpServer.
func NewHTTPServer() *HttpServer {
	return &HttpServer{}
}

// Start starts the HTTP server.
func (s *HttpServer) Start(addr, serviceName string, router http.Handler) {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go s.startServer(srv, serviceName)

	s.waitForShutdown(serviceName, srv)
}

// startServer starts the HTTP server.
func (s *HttpServer) startServer(srv *http.Server, serviceName string) {
	logger.Infof("HTTP server '%s' listening on http://%v", serviceName, srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("HTTP server error: %v", err)
	}
}

// waitForShutdown waits for a shutdown signal to gracefully shut down the server.
func (s *HttpServer) waitForShutdown(serviceName string, srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("Shutting down server '%s'...", serviceName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down server: %v", err)
	}

	logger.Infof("HTTP server '%s' exited.", serviceName)
}
