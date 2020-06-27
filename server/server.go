package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"slowly/logger"
)

type server struct {
	httpServer *http.Server
	logger     logger.Logger
}

type Config struct {
	WriteTimeoutInSec int
	ReadTimeoutInSec  int
	Address           string
}

type Server interface {
	Start()
}

func (s *server) Start() {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan

		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			s.logger.Errorf("Error while shutting down: %s", err.Error())
		}
		close(idleConnsClosed)
	}()

	s.logger.Infof("%v", s.httpServer.ListenAndServe())

	<-idleConnsClosed
}

func NewServer(handler http.Handler, logger logger.Logger, config *Config) Server {
	if config.WriteTimeoutInSec <= 0 {
		panic("WriteTimeoutInSec is not set")
	}
	if config.ReadTimeoutInSec <= 0 {
		panic("ReadTimeoutInSec is not set")
	}
	if len(config.Address) == 0 {
		panic("Address is not set")
	}

	httpserver := http.Server{
		WriteTimeout: time.Duration(config.WriteTimeoutInSec) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeoutInSec) * time.Second,
		Addr:         config.Address,
		Handler:      handler,
	}

	return &server{
		httpServer: &httpserver,
		logger:     logger,
	}
}
