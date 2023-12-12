package main

import (
	"log"

	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/config"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/pkg/db"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/pkg/logger"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/storage"
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.LogLevel)

	db, err := db.New(cfg)
	if err != nil {
		logger.Error("Error while connecting to database", err)
	} else {
		logger.Info("Successfully connected to database")
	}

	router := api.New(logger, cfg, storage.New(db, logger, cfg))

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", err)
	}
}
