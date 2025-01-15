package coinAPI

import (
	"fmt"

	SDK "github.com/CoinAPI/coinapi-sdk/data-api/go-rest/v1"
)

const (
	API_KEY      = "fec05622-dcc5-4a1c-b632-7faf2768ecb6"
	assetIDQoute = "USD"
)

type CoinAPI struct {
	sdk *SDK.SDK
}

func New() *CoinAPI {
	return &CoinAPI{sdk: SDK.NewSDK(API_KEY)}
}

/*func (capi *CoinAPI) WriteExRatesHistory(ctx context.Context, log *slog.Logger, n int) {
	ticker := time.NewTicker(time.Duration(n) * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("ticker ended")
				return
			case <-ticker.C:
				log.Info("rates updated")
			}
		}
	}()
}*/

func (capi *CoinAPI) GetExRate(assetIDBase string) (string, error) {
	const op = "pkg.coinAPI.GetExRate"

	exchange_rat_specific, err := capi.sdk.Exchange_rates_get_specific_rate(assetIDBase, assetIDQoute)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return exchange_rat_specific.Rate.String(), nil
}
