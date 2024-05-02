package main

import (
	"context"
	"github.com/iurikman/smartSurvey/internal/config"
	"github.com/iurikman/smartSurvey/internal/store"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
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
	cfg := config.New()

	err := serverOne.Start(ctx)
	if err != nil {
		log.Panic("Server start error")
	}

	defer zap.L().Sync()

	pgStore, err := store.New(ctx, cfg)
	if err != nil {
		zap.L().With(zap.Error(err)).Panic("pgStore.New")
	}

	if err := pgStore.Migrate(migrate.Up); err != nil {
		zap.L().With(zap.Error(err)).Panic("pgStore.Migrate")
	}

	zap.L().Info("successful migration")
}
