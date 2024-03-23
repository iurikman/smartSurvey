package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	server "github.com/iurikman/smartSurvey/internal/rest"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		<-sigCh
		cancel()
	}()

	serverOne := server.NewServer(":8080")

	err := serverOne.Start(ctx)
	if err != nil {
		log.Panic("Server start error")
	}
}
