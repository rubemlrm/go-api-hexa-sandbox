package api

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
)

func Start(h http.Handler) error {
	server := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stop()
	errShutdown := make(chan error, 1)
	go shutdown(&server, ctx, errShutdown)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	err = <-errShutdown
	if err != nil {
		return err
	}
	return nil
}

func shutdown(server *http.Server, ctxShutdown context.Context, errShutdown chan error) {
	<-ctxShutdown.Done()

	ctxTimeout, stop := context.WithTimeout(context.Background(), 30)
	defer stop()

	err := server.Shutdown(ctxTimeout)
	switch err {
	case nil:
		errShutdown <- nil
	case context.DeadlineExceeded:
		errShutdown <- fmt.Errorf("Forcing closing the server")
	default:
		errShutdown <- fmt.Errorf("Forcing closing the server")
	}
}
