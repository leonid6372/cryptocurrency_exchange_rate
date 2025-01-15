package cryptoExRate

import (
	"context"
	"crypto_ex_rate/internal/storage/postgres"
	"crypto_ex_rate/pkg/coinAPI"
	"crypto_ex_rate/pkg/logger/sl"
	"fmt"
	"log/slog"
	"slices"
	"time"
)

var (
	ErrCryptoIsNotMonitored = "cryptocurrency is not monitored"
)

type CryptoExRate struct {
	log                *slog.Logger
	ctx                context.Context
	secToUpdate        int
	cryptoNames        []string
	coinAPI            *coinAPI.CoinAPI
	cryptoListModifier CryptoListModifier
	cryptoExRateGetter CryptoExRateGetter
}

type CryptoListModifier interface {
	AddCryptoExRate(ctx context.Context, coin, price string) error
}

type CryptoExRateGetter interface {
	GetPrice(ctx context.Context, coin string, timestamp int) (string, error)
}

func New(ctx context.Context, log *slog.Logger, secToUpdate int, coinAPI *coinAPI.CoinAPI, storage *postgres.Storage) *CryptoExRate {
	return &CryptoExRate{
		log:                log,
		ctx:                ctx,
		secToUpdate:        secToUpdate,
		cryptoNames:        []string{},
		coinAPI:            coinAPI,
		cryptoListModifier: storage,
		cryptoExRateGetter: storage,
	}
}

func (cer *CryptoExRate) StartCryptoUpdate() {
	const op = "internal.service.StartCryptoUpdate"
	log := cer.log.With(slog.String("op", op))

	ticker := time.NewTicker(time.Duration(cer.secToUpdate) * time.Second)

	go func() {
		for {
			select {
			case <-cer.ctx.Done():
				log.Info("exchnage rate updating ended")
				return
			case <-ticker.C:
				for _, crypto := range cer.cryptoNames {
					go func(crypto string) {
						price, err := cer.coinAPI.GetExRate(crypto)
						if err != nil {
							log.Error("failed to get exchange rate", sl.Err(err))
						}
						cer.cryptoListModifier.AddCryptoExRate(cer.ctx, crypto, price)
					}(crypto)
				}
			}
		}
	}()
}

func (cer *CryptoExRate) AddCurrency(coin string) error {
	const op = "internal.service.AddCurrency"
	log := cer.log.With(slog.String("op", op))

	price, err := cer.coinAPI.GetExRate(coin)
	if err != nil {
		log.Error("failed to get exchange rate", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = cer.cryptoListModifier.AddCryptoExRate(cer.ctx, coin, price)
	if err != nil {
		log.Error("failed to add cryptocurrency to monitoring", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	cer.cryptoNames = append(cer.cryptoNames, coin)

	return nil
}

func (cer *CryptoExRate) RemoveCurrency(coin string) error {
	const op = "internal.service.RemoveCurrency"
	log := cer.log.With(slog.String("op", op))

	idx := slices.Index(cer.cryptoNames, coin)
	if idx == -1 {
		log.Error("failed to remove cryptocurrency from monitoring")
		return fmt.Errorf("%s: %s", op, ErrCryptoIsNotMonitored)
	}

	cer.cryptoNames[idx] = cer.cryptoNames[len(cer.cryptoNames)-1]
	cer.cryptoNames = cer.cryptoNames[:len(cer.cryptoNames)-1]

	return nil
}

func (cer *CryptoExRate) GetPrice(coin string, timestamp int) (string, error) {
	const op = "internal.service.GetPrice"
	log := cer.log.With(slog.String("op", op))

	price, err := cer.cryptoExRateGetter.GetPrice(cer.ctx, coin, timestamp)
	if err != nil {
		log.Error("failed to get cryptocurrency price", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return price, nil
}
