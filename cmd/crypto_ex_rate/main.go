package main

import (
	"context"
	server "crypto_ex_rate/internal/server/http"
	cryptoExRate "crypto_ex_rate/internal/service"
	"crypto_ex_rate/internal/storage/postgres"
	"crypto_ex_rate/pkg/coinAPI"
	"crypto_ex_rate/pkg/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	secToUpdate   = 2 // exchange rates will be updated every N seconds
	httpSrvAddr   = "0.0.0.0:8085"
	migrationPath = "file://migrations/postgres/"
	postgresInfo  = "host=crypto_ex_rate_db port=5435 user=admin password=1111 dbname=cryptoExRateDB sslmode=disable"
	postgresURL   = "postgres://admin:1111@crypto_ex_rate_db:5435/cryptoExRateDB?sslmode=disable"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, cancel := context.WithCancel(context.Background())

	log.Info("starting application")

	// init and migrate postgres storage
	log.Info("connecting storage")
	storage, err := postgres.New("postgres", postgresInfo)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	err = storage.MigrationUp(postgresURL, migrationPath)
	if err != nil {
		log.Error("failed to migrate storage", sl.Err(err))
		os.Exit(1)
	}

	// init coinAPI
	coinAPI := coinAPI.New()

	// init crypto exchange rate service and update func
	cryptoExRate := cryptoExRate.New(ctx, log, secToUpdate, coinAPI, storage)
	cryptoExRate.StartCryptoUpdate()

	// init http server
	serverAPI := server.New(log, httpSrvAddr, cryptoExRate)
	go func() {
		err = serverAPI.Start(log)
		if err != nil {
			log.Error("failed to init http server", sl.Err(err))
			os.Exit(1)
		}
	}()

	// wait until ctrl+c pressed by user to exit
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()
	log.Info("closing application")

	// gracefull close storage connection
	log.Info("closing storage")
	err = storage.Stop()
	if err != nil {
		log.Error("failed to close storage connection correctly", sl.Err(err))
	}

	// gracefull stop http server
	log.Info("closing http server")
	err = serverAPI.Stop()
	if err != nil {
		log.Error("failed to close http server correctly", sl.Err(err))
	}

	log.Info("application closed")
}
