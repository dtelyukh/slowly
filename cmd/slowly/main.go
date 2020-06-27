package main

import (
	"log"
	"os"
	"runtime/debug"

	"slowly/api"
	"slowly/preferences"
	"slowly/server"

	"github.com/sirupsen/logrus"
)

func handleExit(logger *logrus.Logger) {
	if e := recover(); e != nil {
		logger.Errorf("%v Stask trace: %s", e, debug.Stack())
		os.Exit(1)
	}
}

func main() {
	logger := logrus.New()

	p, err := preferences.Get()
	if err != nil {
		logger.Fatalf("Failed to set preferences: %v\n", err)
	}

	if p.LogAsJSON {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.Level = logrus.Level(p.LogLevel)
	log.SetOutput(logger.Writer())

	defer handleExit(logger)

	handler := api.NewHandler(logger)

	server := server.NewServer(
		handler,
		logger,
		&server.Config{
			WriteTimeoutInSec: p.WriteTimeoutInSec,
			ReadTimeoutInSec:  p.ReadTimeoutInSec,
			Address:           p.HttpAddress,
		})
	server.Start()
}
