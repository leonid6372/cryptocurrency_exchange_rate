package server

import (
	"context"
	"crypto_ex_rate/internal/server/http/handlers/add"
	"crypto_ex_rate/internal/server/http/handlers/price"
	"crypto_ex_rate/internal/server/http/handlers/remove"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type CryptoExRate interface {
	AddCurrency(coin string) error
	RemoveCurrency(coin string) error
	GetPrice(coin string, timestamp int) (string, error)
}

type serverAPI struct {
	srv          *http.Server
	cryptoExRate CryptoExRate
}

func New(log *slog.Logger, address string, cryptoExRate CryptoExRate) serverAPI {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)    // Log all requests
	mux.Use(middleware.Recoverer) // Avoidng panic
	mux.Use(middleware.URLFormat) // URL parser

	mux.Post("/currency/add", add.New(log, cryptoExRate))
	mux.Post("/currency/remove", remove.New(log, cryptoExRate))
	mux.Post("/currency/price", price.New(log, cryptoExRate))

	server := http.Server{Addr: address, Handler: mux}
	return serverAPI{srv: &server, cryptoExRate: cryptoExRate}
}

func (srvAPI *serverAPI) Start(log *slog.Logger) error {
	const op = "internal.server.http.Start"

	log.Info("starting http server: " + srvAPI.srv.Addr)
	err := srvAPI.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (srvAPI *serverAPI) Stop() error {
	const op = "internal.server.http.Stop"

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srvAPI.srv.Shutdown(shutDownCtx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
