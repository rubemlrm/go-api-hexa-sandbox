package api

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/pkg/gin"
)

type Server struct {
	*http.Server
}

func NewServer(h http.Handler, httpConfig config.HTTP, logger *slog.Logger) (Server, error) {
	server, err := prepareServerConfig(h, httpConfig, logger)
	if err != nil {
		return Server{}, err
	}
	return server, err
}

func (s Server) Start() error {
	// init empty variable to store errors from go routines
	var err error

	// Set channel that will receive signals
	notifyChannel := make(chan os.Signal, 1)
	signal.Notify(notifyChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start Server list
	go func() {
		err = s.startListening()
	}()

	// Wait for signal to shutdown
	shutdownSignal := <-notifyChannel
	log.Printf("Received signal: %v\n", shutdownSignal)

	go func() {
		err = s.stopListening()
	}()

	println(fmt.Sprintf("Current service listening on port %s \n", s.Addr))

	return err
}

func prepareServerConfig(h http.Handler, httpConfig config.HTTP, logger *slog.Logger) (Server, error) {
	_, err := strconv.Atoi(httpConfig.ReadTimeout)

	if err != nil {
		return Server{}, &gin.HTTPConfigurationError{Input: "ReadTimeout"}
	}

	_, err = strconv.Atoi(httpConfig.WriteTimeout)

	if err != nil {
		return Server{}, &gin.HTTPConfigurationError{Input: "WriteTimeout"}
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", httpConfig.Address),
		Handler:           h,
		ReadHeaderTimeout: 120,
	}
	logger.Info("server", "configuration", server)
	return Server{server}, nil
}

func (s *Server) startListening() error {
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) stopListening() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
