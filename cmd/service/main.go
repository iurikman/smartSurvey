package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/iurikman/smartSurvey/internal/config"
	"github.com/iurikman/smartSurvey/internal/logger"
	server "github.com/iurikman/smartSurvey/internal/rest"
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
	serverOne := server.NewServer(server.Config{
		BindAddress: cfg.BindAddress,
	})

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

	log.Info("successful migration")

	err = serverOne.Start(ctx)
	if err != nil {
		log.Panicf("Server start error: %v", err)
	}
	serverOne.CreateUser("test user")
}
