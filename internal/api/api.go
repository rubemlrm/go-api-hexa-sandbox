package api

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rubemlrm/go-api-bootstrap/config"
)

func Start(h http.Handler, httpConfig config.HTTP) error {

	server, err := prepareServerConfig(h, httpConfig)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stop()
	errShutdown := make(chan error, 1)
	go shutdown(server, ctx, errShutdown)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	err = <-errShutdown
	if err != nil {
		return err
	}
	return nil
}

func prepareServerConfig(h http.Handler, httpConfig config.HTTP) (*http.Server, error) {

	_, err := strconv.Atoi(httpConfig.ReadTimeout)

	if err != nil {
		return nil, err
	}

	_, err = strconv.Atoi(httpConfig.WriteTimeout)

	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    ":" + httpConfig.Address,
		Handler: h,
	}

	return server, nil

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
