package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/iurikman/smartSurvey/internal/config"
	"github.com/iurikman/smartSurvey/internal/logger"
	server "github.com/iurikman/smartSurvey/internal/rest"
	"github.com/iurikman/smartSurvey/internal/service"
	"github.com/iurikman/smartSurvey/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.InitLogger("debug")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	defer cancel()

	cfg := config.New()

	pgStore, err := store.New(ctx, store.Config{
		PGUser:     cfg.PGUser,
		PGPassword: cfg.PGPassword,
		PGHost:     cfg.PGHost,
		PGPort:     cfg.PGPort,
		PGDatabase: cfg.PGDatabase,
	})
	if err != nil {
		log.Panicf("pgStore.New: %v", err)
	}

	if err := pgStore.Migrate(migrate.Up); err != nil {
		log.Panicf("pgStore.Migrate: %v", err)
	}

	svc := service.New(pgStore)

	httpServer := server.NewServer(
		server.Config{BindAddress: cfg.BindAddress},
		svc,
	)

	log.Info("successful migration")

	if err = httpServer.Start(ctx); err != nil {
		log.Panicf("Server start error: %v", err)
	}
}
