package price

import (
	resp "crypto_ex_rate/pkg/api/response"
	"crypto_ex_rate/pkg/logger/sl"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Coin      string `json:"coin" validate:"required"`
	Timestamp int    `json:"timestamp" validate:"required"`
}

type Response struct {
	resp.Response
	ExchangeRate string `json:"exchange_rate"`
}

type CryptoExRate interface {
	GetPrice(coin string, timestamp int) (string, error)
}

func New(log *slog.Logger, cryptoExRate CryptoExRate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http.handlers.add.New"
		log := log.With(slog.String("op", op))

		var req Request

		// Decode data from request JSON
		err := render.DecodeJSON(r.Body, &req)
		// Exception: empty request
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			w.WriteHeader(400)
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(400)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Validate required request fields
		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(400)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		price, err := cryptoExRate.GetPrice(req.Coin, req.Timestamp)
		if err != nil {
			log.Error("failed to get cryptocurrency exchange rate", sl.Err(err))
			w.WriteHeader(500)
			render.JSON(w, r, resp.Error("failed to get cryptocurrency exchange rate"))
			return
		}

		// Marshal data to correct JSON response
		response, err := json.Marshal(Response{Response: resp.OK(), ExchangeRate: price})
		if err != nil {
			log.Error("failed to process response", sl.Err(err))
			w.WriteHeader(500)
			render.JSON(w, r, resp.Error("failed to process response"))
			return
		}
		render.Data(w, r, response)
	}
}
