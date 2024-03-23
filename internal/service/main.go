package main

import (
	"context"
	"github.com/iurikman/smartSurvey/internal"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var log = logrus.New()

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
