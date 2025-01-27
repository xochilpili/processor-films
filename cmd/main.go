package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/logger"
	"github.com/xochilpili/processor-films/internal/webserver"
)

func main() {
	config := config.New()
	logger := logger.New()

	srv := webserver.New(config, logger)
	go func() {
		logger.Info().Msgf("starting server at %s:%s", config.Host, config.Port)
		if err := srv.Web.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("error while loading server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-shutdown
	logger.Info().Msg("shutting down server")
	if err := srv.Web.Shutdown(context.Background()); err != nil {
		logger.Fatal().Err(err).Msg("error while shutting down server.")
	}

}
