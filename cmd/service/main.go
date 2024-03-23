package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/iurikman/smartSurvey/internal/logger"
	server "github.com/iurikman/smartSurvey/internal/rest"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.InitLogger("debug")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	defer cancel()

	serverOne := server.NewServer(":8080")

	err := serverOne.Start(ctx)
	if err != nil {
		log.Panic("Server start error")
	}
}
