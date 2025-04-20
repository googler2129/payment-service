package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mercor/payment-service/pkg/env"
	"github.com/mercor/payment-service/pkg/log"
	"github.com/mercor/payment-service/pkg/shutdown"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	*http.Server
}

// InitializeServer InitializeRouter Returns a new http server which internally uses gin frameworks. It registers all endpoints and middlewares.
func InitializeServer(
	listenAddr string,
	readTimeout,
	writeTimeout,
	idleTimeout time.Duration,
	isRedirectTrailingSlash bool,
	middlewares ...gin.HandlerFunc,
) *Server {

	//Recovery middleware recovers from any panics and writes a 500 if there was one
	//Request ID middleware adds request ID in every request if not present in header
	defaultMiddlewares := []gin.HandlerFunc{gin.Recovery(), env.RequestID()}

	// Setting gin to releaseMode
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Set RedirectTrailingSlash to false
	r.RedirectTrailingSlash = isRedirectTrailingSlash

	for _, middleware := range defaultMiddlewares {
		r.Use(middleware)
	}

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	if readTimeout == 0 {
		readTimeout = 10 * time.Second
	}

	if writeTimeout == 0 {
		writeTimeout = 10 * time.Second
	}

	if idleTimeout == 0 {
		idleTimeout = 70 * time.Second
	}

	s := &http.Server{
		Addr:         listenAddr,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	server := &Server{
		r,
		s,
	}

	return server
}

// StartServer Starts the http server
func (s *Server) StartServer(serviceName string) error {
	// Register shutdown callback
	shutdown.RegisterShutdownCallback(serviceName, s)

	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context Timeout:
		log.WithError(err).Panic("HTTP server Shutdown error: %v", err)
		return err
	}
	return nil
}
